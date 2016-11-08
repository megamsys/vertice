package machine

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	nsqp "github.com/crackcomm/nsqueue/producer"
	"github.com/megamsys/libgo/events"
	"github.com/megamsys/libgo/events/alerts"
	"github.com/megamsys/libgo/utils"
	constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/opennebula-go/compute"
	"github.com/megamsys/opennebula-go/disk"
	"github.com/megamsys/opennebula-go/images"
	"github.com/megamsys/opennebula-go/virtualmachine"
	"github.com/megamsys/vertice/carton"
	lb "github.com/megamsys/vertice/logbox"
	"github.com/megamsys/vertice/meta"
	"github.com/megamsys/vertice/provision"
	"github.com/megamsys/vertice/provision/one/cluster"
)

type OneProvisioner interface {
	Cluster() *cluster.Cluster
}

type Machine struct {
	Name         string
	Region       string
	Id           string
	CartonId     string
	CartonsId    string
	AccountsId   string
	Level        provision.BoxLevel
	SSH          provision.BoxSSH
	Image        string
	VCPUThrottle string
	VMId         string
	VNCHost      string
	VNCPort      string
	ImageId      string
	Routable     bool
	Status       utils.Status
	State        utils.State
}

type CreateArgs struct {
	Commands    []string
	Box         *provision.Box
	Compute     provision.BoxCompute
	Deploy      bool
	Provisioner OneProvisioner
}

func (m *Machine) Create(args *CreateArgs) error {
	opts := compute.VirtualMachine{
		Name:   m.Name,
		Image:  m.Image,
		Region: args.Box.Region,
		Cpu:    strconv.FormatInt(int64(args.Box.GetCpushare()), 10),
		Memory: strconv.FormatInt(int64(args.Box.GetMemory()), 10),
		HDD:    strconv.FormatInt(int64(args.Box.GetHDD()), 10),
		ContextMap: map[string]string{compute.ASSEMBLY_ID: args.Box.CartonId,compute.ORG_ID: args.Box.OrgId,
			compute.ASSEMBLIES_ID: args.Box.CartonsId, compute.ACCOUNTS_ID: args.Box.AccountsId},
		Vnets: args.Box.Vnets,
	}

	//m.addEnvsToContext(m.BoxEnvs, &vm)
	_, _, vmid, err := args.Provisioner.Cluster().CreateVM(opts, m.VCPUThrottle)
	if err != nil {
		return err
	}
	m.VMId = vmid

	var id = make(map[string][]string)
	vm := []string{}
	vm = []string{m.VMId}
	id[carton.VMID] = vm
	if asm, err := carton.NewAmbly(m.CartonId); err != nil {
		return err
	} else if err = asm.NukeAndSetOutputs(id); err != nil {
		return err
	}
	return nil
}

func (m *Machine) VmHostIpPort(args *CreateArgs) error {

	opts := virtualmachine.Vnc{
		VmId: m.VMId,
	}

	vnchost, vncport, err := args.Provisioner.Cluster().GetIpPort(opts, m.Region)
	if err != nil {
		return err
	}
	m.VNCHost = vnchost
	m.VNCPort = vncport
	return nil
}


func (m *Machine) UpdateVncHostPost() error {
	var vnc = make(map[string][]string)
	var port, host  []string
	host = []string{m.VNCHost}
	port = []string{m.VNCPort}
	vnc[carton.VNCHOST] = host
	vnc[carton.VNCPORT] = port
	if asm, err := carton.NewAmbly(m.CartonId); err != nil {
		return err
	} else if err = asm.NukeAndSetOutputs(vnc); err != nil {
		return err
	}
	return nil
}

func (m *Machine) Remove(p OneProvisioner) error {
	log.Debugf("  removing machine in one (%s)", m.Name)
	id, _ := strconv.Atoi(m.VMId)
	opts := compute.VirtualMachine{
		Name:   m.Name,
		Region: m.Region,
		VMId:   id,
	}

	err := p.Cluster().DestroyVM(opts)
	if err != nil {
		return err
	}
	return nil
}

//trigger multi event in the order
func (m *Machine) Deduct() error {
	mi := make(map[string]string)
	mi[constants.ACCOUNTID] = m.AccountsId
	mi[constants.ASSEMBLYID] = m.CartonId
	mi[constants.ASSEMBLYNAME] = m.Name
	mi[constants.CONSUMED] = "0.1"
	mi[constants.START_TIME] = time.Now().Add(-10 * time.Minute).String()
	mi[constants.END_TIME] = time.Now().String()
	newEvent := events.NewMulti(
		[]*events.Event{
			&events.Event{
				AccountsId:  m.AccountsId,
				EventAction: alerts.DEDUCT,
				EventType:   constants.EventBill,
				EventData:   alerts.EventData{M: mi},
				Timestamp:   time.Now().Local(),
			},
			&events.Event{
				AccountsId:  m.AccountsId,
				EventAction: alerts.TRANSACTION, //Change type to transaction
				EventType:   constants.EventBill,
				EventData:   alerts.EventData{M: mi},
				Timestamp:   time.Now().Local(),
			},
		})
	return newEvent.Write()
}

func (m *Machine) LifecycleOps(p OneProvisioner, action string) error {
	log.Debugf("  %s machine in one (%s)", action, m.Name)
	id, _ := strconv.Atoi(m.VMId)
	opts := compute.VirtualMachine{
		Name:   m.Name,
		Region: m.Region,
		VMId:   id,
	}
	err := p.Cluster().VM(opts, action)
	if err != nil {
		return err
	}
	return nil
}

//it possible to have a Notifier interface that does this, duck typed b y Assembly, Components.
func (m *Machine) SetStatus(status utils.Status) error {
	log.Debugf("  set status[%s] of machine (%s, %s)", m.Id, m.Name, status.String())

	if asm, err := carton.NewAmbly(m.CartonId); err != nil {
		return err
	} else if err = asm.SetStatus(status); err != nil {
		return err
	}

	if m.Level == provision.BoxSome {
		log.Debugf("  set status[%s] of machine (%s, %s)", m.Id, m.Name, status.String())

		if comp, err := carton.NewComponent(m.Id); err != nil {
			return err
		} else if err = comp.SetStatus(status); err != nil {
			return err
		}
	}
	return nil
}


func (m *Machine) SetMileStone(state utils.State) error {
	log.Debugf("  set state[%s] of machine (%s, %s)", m.Id, m.Name, state.String())

	if asm, err := carton.NewAmbly(m.CartonId); err != nil {
		return err
	} else if err = asm.SetState(state); err != nil {
		return err
	}

	if m.Level == provision.BoxSome {
		log.Debugf("  set state[%s] of machine (%s, %s)", m.Id, m.Name, state.String())

		if comp, err := carton.NewComponent(m.Id); err != nil {
			return err
		} else if err = comp.SetState(state); err != nil {
			return err
		}
	}
	return nil
}

//just publish a message stateup to the machine.
func (m *Machine) ChangeState(status utils.Status) error {
	log.Debugf("  change state of machine (%s, %s)", m.Name, status.String())

	pons := nsqp.New()
	if err := pons.Connect(meta.MC.NSQd[0]); err != nil {
		return err
	}

	bytes, err := json.Marshal(
		carton.Requests{
			CatId:     m.CartonId,
			Action:    status.String(),
			Category:  carton.STATE,
			CreatedAt: time.Now().String(),
		})

	if err != nil {
		return err
	}

	log.Debugf("  pub to machine (%s, %s)", m.Name, bytes)

	if err = pons.Publish(m.Name, bytes); err != nil {
		return err
	}

	defer pons.Stop()
	return nil
}

//if there is a file or something to be created, do it here.
func (m *Machine) Logs(p OneProvisioner, w io.Writer) error {

	fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("logs nirvana ! machine %s ", m.Name)))
	return nil
}

func (m *Machine) Exec(p OneProvisioner, stdout, stderr io.Writer, cmd string, args ...string) error {
	cmds := []string{"/bin/bash", "-lc", cmd}
	cmds = append(cmds, args...)

	//load the ssh key inmemory
	//ssh and run the command
	//sshOpts := ssh.CreateExecOptions{
	//}

	//if err != nil {
	//	return err
	//}

	return nil

}

func (m *Machine) SetRoutable(ip string) {
	m.Routable = (len(strings.TrimSpace(ip)) > 0)
}

func (m *Machine) addEnvsToContext(envs string, cfg *compute.VirtualMachine) {
	/*
		for _, envData := range envs {
			cfg.Env = append(cfg.Env, fmt.Sprintf("%s=%s", envData.Name, envData.Value))
		}
			cfg.Env = append(cfg.Env, []string{
				fmt.Sprintf("%s=%s", "MEGAM_HOST", host),
			}...)
	*/
}

func (m *Machine) CreateDiskSnap(p OneProvisioner) error {
		log.Debugf("  creating snap machine in one (%s)", m.Name)
	snp, err := carton.GetSnap(m.CartonsId)
	if err != nil {
		return err
	}

  vmid,_ := strconv.Atoi(m.VMId)
	opts := compute.Image{
		Name:   snp.Name,
		Region: m.Region,
		VMId:   vmid,
	}

	id, err := p.Cluster().SnapVMDisk(opts)
	if err != nil {
		return err
	}
  m.ImageId = id
	return nil
}

func (m * Machine) IsSnapReady(p OneProvisioner) error {
  id, _ := strconv.Atoi(m.ImageId)
	opts := &images.Image{
		Id: id,
	}
	err := p.Cluster().IsSnapReady(opts,m.Region)
  if err != nil {
	  return err
  }

return nil
}

//it possible to have a Notifier interface that does this, duck typed by Snap id
func (m *Machine) AttachNewDisk(p OneProvisioner) error {

	dsk, err := carton.GetDisks(m.CartonsId)
	if err != nil {
		return err
	}
  size := dsk.NumMemory()
  id,_ := strconv.Atoi(m.VMId)
	opts := &disk.VmDisk{
		VmId: id,
		Vm:  disk.Vm{Disk: disk.Disk{Size: size}},
	}

	err = p.Cluster().AttachDisk(opts,m.Region)
	if err != nil {
		return err
	}
	return nil
}

func (m *Machine) UpdateSnap() error {
	update_fields := make(map[string]interface{},2)
	update_fields["Image_Id"] = m.ImageId
	update_fields["Status"] = "ready"
  sns,err := carton.GetSnap(m.CartonsId)
	if err != nil {
		return err
	}
	//d := &carton.Snaps{Id: m.CartonsId}

	err = sns.UpdateSnap(update_fields)
	if err != nil {
		return err
	}
	return nil
}

func (m *Machine) UpdateDisk(p OneProvisioner) error {
	id ,_ := strconv.Atoi(m.VMId)
  vd := &disk.VmDisk{VmId: id}

	l, err := p.Cluster().GetDiskId(vd,m.Region)
	if err != nil {
		return err
	}
	update_fields := make(map[string]interface{},2)
	update_fields["Disk_Id"] = strconv.Itoa(l[len(l)-1])
	update_fields["Status"] = "Success"

	d,err := carton.GetDisks(m.CartonsId)
	if err != nil {
		return err
	}
	err = d.UpdateDisk(update_fields)
	if err != nil {
		return err
	}
	return nil
}

func (m *Machine) RemoveDisk(p OneProvisioner) error {
	dsk, err := carton.GetDisks(m.CartonsId)
	if err != nil {
		return err
	}

  id,_ := strconv.Atoi(m.VMId)
	did,_ := strconv.Atoi(dsk.DiskId)
	opts := &disk.VmDisk{
		VmId:  id,
		Vm:   disk.Vm{Disk: disk.Disk{Disk_Id: did}},
	}

	err = p.Cluster().DetachDisk(opts,m.Region)
	if err != nil {
		return err
	}

	err = dsk.RemoveDisk()
	if err != nil {
		return err
	}

	return nil
}

func (m *Machine) RemoveSnapshot(p OneProvisioner) error {
	snp, err := carton.GetSnap(m.CartonsId)
	if err != nil {
		return err
	}
  id, _ := strconv.Atoi(snp.ImageId)
	log.Debugf("  remove snap machine in one (%s)", m.Name)
	opts := compute.Image{
		Name:   snp.Name,
		Region: m.Region,
		ImageId: id,
	}
	err = p.Cluster().RemoveSnap(opts)
	if err != nil {
		return err
	}

	err = snp.RemoveSnap()
	if err != nil {
		return err
	}

	return nil
}
