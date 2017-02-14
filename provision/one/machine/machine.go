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
	"github.com/megamsys/libgo/events/alerts"
	"github.com/megamsys/libgo/events/bills"
	"github.com/megamsys/libgo/safe"
	"github.com/megamsys/libgo/utils"
	constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/opennebula-go/compute"
	"github.com/megamsys/opennebula-go/snapshot"
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
	AccountId    string
	Level        provision.BoxLevel
	SSH          provision.BoxSSH
	Image        string
	VCPUThrottle string
	VMId         string
	VNCHost      string
	VNCPort      string
	ImageId      string
	StorageType  string
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
	asm, err := carton.NewAssembly(m.CartonId, m.AccountId, "")
	if err != nil {
		return err
	}

	opts := compute.VirtualMachine{
		Name:       m.Name,
		Image:      m.Image,
		Region:     args.Box.Region,
		Cpu:        strconv.FormatInt(int64(args.Box.GetCpushare()), 10),
		Memory:     strconv.FormatInt(int64(args.Box.GetMemory()), 10),
		HDD:        strconv.FormatInt(int64(args.Box.GetHDD()), 10),
		CpuCost:    asm.GetVMCpuCost(),
		MemoryCost: asm.GetVMMemoryCost(),
		HDDCost:    asm.GetVMHDDCost(),
		ContextMap: map[string]string{compute.ASSEMBLY_ID: args.Box.CartonId, compute.ORG_ID: args.Box.OrgId,
			compute.ASSEMBLIES_ID: args.Box.CartonsId, compute.ACCOUNTS_ID: args.Box.AccountId, compute.API_KEY: args.Box.ApiArgs.Api_Key, constants.QUOTA_ID: args.Box.QuotaId},
		Vnets: args.Box.Vnets,
	}
	opts.VCpu = opts.Cpu
	if strings.Contains(args.Box.Tosca, "freebsd") {
		opts.Files = "/detio/freebsd/init.sh"
	}

	_, _, vmid, err := args.Provisioner.Cluster().CreateVM(opts, m.VCPUThrottle, m.StorageType)
	if err != nil {
		return err
	}
	m.VMId = vmid

	var id = make(map[string][]string)
	vm := []string{}
	vm = []string{m.VMId}
	id[carton.INSTANCE_ID] = vm

	if err = asm.NukeAndSetOutputs(id); err != nil {
		return err
	}
	return nil
}

func (m *Machine) CheckCredits(b *provision.Box, w io.Writer) error {
	bal, err := bills.NewBalances(b.AccountId, meta.MC.ToMap())
	if err != nil || bal == nil {
		return err
	}

	//have to decide what to do whether balance record is empty
	i, err := strconv.ParseFloat(bal.Credit, 64)
	if err != nil {
		return err
	}

	if i <= 0 {
		carton.DoneNotify(b, w, alerts.INSUFFICIENT_FUND)
		log.Debugf(" credit balance insufficient for the user (%s)", b.AccountId)
		return fmt.Errorf("credit balance insufficient")
	}

	return nil
}

func (m *Machine) VmHostIpPort(args *CreateArgs) error {
	asm, err := carton.NewAssembly(m.CartonId, m.AccountId, "")
	if err != nil {
		return err
	}
	opts := virtualmachine.Vnc{
		VmId: m.VMId,
	}

	res := &virtualmachine.VM{}
	_ = asm.SetStatus(utils.Status(constants.StatusLcmStateChecking))

	err = safe.WaitCondition(60*time.Minute, 20*time.Second, func() (bool, error) {
		_ = asm.Trigger_event(utils.Status(constants.StatusWaitUntill))
		res, err = args.Provisioner.Cluster().GetVM(opts, m.Region)
		if err != nil {
			return false, err
		}
		status := res.StateString()
		if res.LcmStateString() != "" {
			status = status + "_" + res.LcmStateString()
		}
		_ = asm.Trigger_event(utils.Status(status))
		return (res.HistoryRecords.History != nil && res.LcmState == 3), nil
	})

	if err != nil {
		return err
	}

	m.VNCHost = res.GetHostIp()
	m.VNCPort = res.GetPort()
	return nil
}

func (m *Machine) WaitUntillVMState(args *CreateArgs, vm virtualmachine.VmState, lcm virtualmachine.LcmState) error {
	opts := virtualmachine.Vnc{VmId: m.VMId}

	err := safe.WaitCondition(10*time.Minute, 15*time.Second, func() (bool, error) {
		res, err := args.Provisioner.Cluster().GetVM(opts, m.Region)
		if err != nil {
			return false, err
		}
		if res.IsFailure() {
			return false, fmt.Errorf(res.UserTemplate.Error)
		}
		return (res.State == int(vm) && res.LcmState == int(lcm)), nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (m *Machine) UpdateVncHostPost() error {
	var vnc = make(map[string][]string)
	var port, host []string
	host = []string{m.VNCHost}
	port = []string{m.VNCPort}
	vnc[carton.VNCHOST] = host
	vnc[carton.VNCPORT] = port
	if asm, err := carton.NewAssembly(m.CartonId, m.AccountId, ""); err != nil {
		return err
	} else if err = asm.NukeAndSetOutputs(vnc); err != nil {
		return err
	}
	return nil
}

func (m *Machine) UpdateVMIps(p OneProvisioner) error {
	opts := virtualmachine.Vnc{
		VmId: m.VMId,
	}
	res, err := p.Cluster().GetVM(opts, m.Region)
	if err != nil {
		return err
	}
	ips := m.mergeSameIPtype(m.IPs(res.Nics()))
	log.Debugf("  find and setips of machine (%s, %s)", m.Id, m.Name)
	asm, err := carton.NewAssembly(m.CartonId, m.AccountId, "")
	if  err != nil {
		return err
	}
	return asm.NukeAndSetOutputs(ips)
}

func (m *Machine) IPs(nics []virtualmachine.Nic) map[string][]string {
	var ips = make(map[string][]string)
	pubipv4s := []string{}
	priipv4s := []string{}
	for _, nic := range nics {
			if nic.IPaddress != "" {
				ip4 := strings.Split(nic.IPaddress, ".")
				if len(ip4) == 4 {
					if (ip4[0] == "192" || ip4[0] == "10" || ip4[0] == "172") {
						priipv4s = append(priipv4s, nic.IPaddress)
					} else {
						pubipv4s = append(pubipv4s, nic.IPaddress)
					}
				}
			}
	}

ips[carton.PUBLICIPV4] = pubipv4s
ips[carton.PRIVATEIPV4] = priipv4s
return ips
}

func (m *Machine) mergeSameIPtype(mm map[string][]string)  map[string][]string {
  for IPtype, ips := range mm {
		var sameIp string
		for _, ip := range ips {
			sameIp = sameIp +  ip + ", "
		}
		if sameIp != "" {
			mm[IPtype] = []string{strings.TrimRight(sameIp, ", ")}
		}
	}
	return mm
}


func (m *Machine) Remove(p OneProvisioner, state constants.State) error {
	log.Debugf("  removing machine in one (%s)", m.Name)
	id, _ := strconv.Atoi(m.VMId)
	opts := compute.VirtualMachine{
		Name:   m.Name,
		Region: m.Region,
		VMId:   id,
	}

	err := p.Cluster().ForceDestoryVM(opts)
	if err != nil {
		return err
	}
	return nil
}

func isDeleteOk(state constants.State) bool {
	return state != constants.StateInitialized && state != constants.StateInitializing && state != constants.StatePreError
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

	if asm, err := carton.NewAssembly(m.CartonId, m.AccountId, ""); err != nil {
		return err
	} else if err = asm.SetStatus(status); err != nil {
		return err
	}

	if m.Level == provision.BoxSome {
		log.Debugf("  set status[%s] of machine (%s, %s)", m.Id, m.Name, status.String())

		if comp, err := carton.NewComponent(m.Id, m.AccountId, ""); err != nil {
			return err
		} else if err = comp.SetStatus(status, m.AccountId); err != nil {
			return err
		}
	}
	return nil
}

func (m *Machine) SetMileStone(state utils.State) error {
	log.Debugf("  set state[%s] of machine (%s, %s)", m.Id, m.Name, state.String())

	if asm, err := carton.NewAssembly(m.CartonId, m.AccountId, ""); err != nil {
		return err
	} else if err = asm.SetState(state); err != nil {
		return err
	}

	if m.Level == provision.BoxSome {
		log.Debugf("  set state[%s] of machine (%s, %s)", m.Id, m.Name, state.String())

		if comp, err := carton.NewComponent(m.Id, m.AccountId, ""); err != nil {
			return err
		} else if err = comp.SetState(state, m.AccountId); err != nil {
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
			AccountId: m.AccountId,
			CreatedAt: time.Now(),
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

func (m *Machine) CreateDiskImage(p OneProvisioner) error {
	log.Debugf("  creating snap machine in one (%s)", m.Name)
	bk, err := carton.GetBackup(m.CartonsId, m.AccountId)
	if err != nil {
		return err
	}

	vmid, _ := strconv.Atoi(m.VMId)
	opts := compute.Image{
		Name:   bk.Name,
		Region: m.Region,
		VMId:   vmid,
		DiskId: 0,
		SnapId: -1,
	}

	id, err := p.Cluster().SaveDiskImage(opts)
	if err != nil {
		return err
	}
	m.ImageId = id
	return nil
}


func (m *Machine) CreateDiskSnap(p OneProvisioner) error {
	log.Debugf("  creating snap machine in one (%s)", m.Name)
	snp, err := carton.GetSnap(m.CartonsId, m.AccountId)
	if err != nil {
		return err
	}

	vmid, _ := strconv.Atoi(m.VMId)
	opts := snapshot.Snapshot{
	  VMId: vmid,
	  DiskId: 0,
	  DiskDiscription: snp.Name,
	}

	id, err := p.Cluster().SnapVMDisk(opts,m.Region)
	if err != nil {
		return err
	}
	m.ImageId = id
	return nil
}

func (m *Machine) RestoreSnapshot(p OneProvisioner) error {
	log.Debugf("  restoring snap machine in one (%s)", m.Name)
	snp, err := carton.GetSnap(m.CartonsId, m.AccountId)
	if err != nil {
		return err
	}

	vmid, _ := strconv.Atoi(m.VMId)
	diskId, _ := strconv.Atoi(snp.DiskId)
	sid, _ := strconv.Atoi(snp.SnapId)
	opts := snapshot.Snapshot{
		VMId: vmid,
		DiskId: diskId,
		SnapId: sid,
	}
  m.ImageId = snp.SnapId
	return p.Cluster().RestoreSnap(opts,m.Region)
}

func (m *Machine) IsImageReady(p OneProvisioner) error {
	id, _ := strconv.Atoi(m.ImageId)
	opts := &images.Image{
		Id: id,
	}
	return p.Cluster().IsImageReady(opts, m.Region)
}

func (m *Machine) IsSnapReady(p OneProvisioner) error {
	opts := virtualmachine.Vnc{VmId: m.VMId}
	err := safe.WaitCondition(10*time.Minute, 15*time.Second, func() (bool, error) {
		res, err := p.Cluster().GetVM(opts, m.Region)
		if err != nil {
			return false, err
		}
		if res.IsFailure() {
			return false, fmt.Errorf(res.UserTemplate.Error)
		}
		return res.IsSnapshotReady(), nil
	})

	if err != nil {
		return err
	}
	return nil
}


//it possible to have a Notifier interface that does this, duck typed by Snap id
func (m *Machine) AttachNewDisk(p OneProvisioner) error {
	log.Debugf("  attachng new disk for the machine (%s)", m.Name)
	dsk, err := carton.GetDisks(m.CartonsId, m.AccountId)
	if err != nil {
		return err
	}
	size := dsk.NumMemory()
	id, _ := strconv.Atoi(m.VMId)
	opts := &disk.VmDisk{
		VmId: id,
		Vm:   disk.Vm{Disk: disk.Disk{Size: size}},
	}

	return p.Cluster().AttachDisk(opts, m.Region)
}

func (m *Machine) UpdateSnap() error {
	sns, err := carton.GetSnap(m.CartonsId, m.AccountId)
	if err != nil {
		return err
	}
	sns.SnapId = m.ImageId
	sns.DiskId = "0"
	sns.Status = "created"
	return sns.UpdateSnap()
}

func (m *Machine) MakeActiveSnapshot() error {
	snaps, err := carton.GetAsmSnaps(m.CartonId, m.AccountId)
	if err != nil {
		return err
	}
	for _, v := range snaps {
		if v.SnapId != m.ImageId  {  // && v.Status == constants.ACTIVESNAP
			v.Status = constants.DEACTIVESNAP
			err = v.UpdateSnap()
			if err != nil {
				return err
			}
		} else if v.Id == m.CartonsId {
			v.Status = constants.ACTIVESNAP
			err = v.UpdateSnap()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *Machine) UpdateSnapStatus(status utils.Status) error {
	sns, err := carton.GetSnap(m.CartonsId, m.AccountId)
	if err != nil {
		return err
	}
	sns.Status = status.String()
	return sns.UpdateSnap()
}

func (m *Machine) UpdateBackup() error {
	bk, err := carton.GetBackup(m.CartonsId, m.AccountId)
	if err != nil {
		return err
	}
	bk.ImageId = m.ImageId
	bk.Status = "ready"
	return bk.UpdateBackup()
}

func (m *Machine) UpdateBackupStatus(status utils.Status) error {
	sns, err := carton.GetBackup(m.CartonsId, m.AccountId)
	if err != nil {
		return err
	}
	sns.Status = status.String()
	return sns.UpdateBackup()
}



func (m *Machine) UpdateDisk(p OneProvisioner) error {
	id, _ := strconv.Atoi(m.VMId)
	vd := &disk.VmDisk{VmId: id}

	l, err := p.Cluster().GetDiskId(vd, m.Region)
	if err != nil {
		return err
	}

	d, err := carton.GetDisks(m.CartonsId, m.AccountId)
	if err != nil {
		return err
	}

	d.DiskId = strconv.Itoa(l[len(l)-1])
	d.Status = "success"
	return d.UpdateDisk()
}

func (m *Machine) RemoveDisk(p OneProvisioner) error {
	dsk, err := carton.GetDisks(m.CartonsId, m.AccountId)
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(m.VMId)
	did, _ := strconv.Atoi(dsk.DiskId)
	opts := &disk.VmDisk{
		VmId: id,
		Vm:   disk.Vm{Disk: disk.Disk{Disk_Id: did}},
	}

	return p.Cluster().DetachDisk(opts, m.Region)
}

func (m *Machine) RemoveBackupImage(p OneProvisioner) error {
	bk, err := carton.GetBackup(m.CartonsId, m.AccountId)
	if err != nil {
		return err
	}
	id, _ := strconv.Atoi(bk.ImageId)
	log.Debugf("  remove snap machine in one (%s)", m.Name)
	opts := compute.Image{
		Name:    bk.Name,
		Region:  m.Region,
		ImageId: id,
	}
	err = p.Cluster().RemoveBackup(opts)
	if err != nil {
		return err
	}

	return bk.RemoveBackup()
}


func (m *Machine) RemoveSnapshot(p OneProvisioner) error {
	snp, err := carton.GetSnap(m.CartonsId, m.AccountId)
	if err != nil {
		return err
	}

	vmid, _ := strconv.Atoi(m.VMId)
	diskId, _ := strconv.Atoi(snp.DiskId)
	sid, _ := strconv.Atoi(snp.SnapId)
	opts := snapshot.Snapshot{
		VMId: vmid,
		DiskId: diskId,
		SnapId: sid,
	}

	err = p.Cluster().RemoveSnap(opts, m.Region)
	if err != nil {
		return err
	}

	return snp.RemoveSnap()
}

func (m *Machine) UpdateSnapQuotas(id string) error {
	quota, err := carton.NewQuota(m.AccountId, id)
	if err != nil {
		return err
	}
	count, _ := strconv.Atoi(quota.AllowedSnaps())
	mm := make(map[string][]string, 1)
	if m.Status == constants.StatusSnapCreated {
			mm["no_of_units"] = []string{strconv.Itoa(count - 1)}
	} else if m.Status == constants.StatusSnapDeleted {
			mm["no_of_units"] = []string{strconv.Itoa(count + 1)}
	}
	quota.Allowed.NukeAndSet(mm) //just nuke the matching key:
	return quota.Update()
}

func (m *Machine) IsSnapCreated() bool {
	return m.Status == constants.StatusSnapCreated
}

func (m *Machine) IsSnapDeleted() bool {
	return m.Status == constants.StatusSnapDeleted
}

func (m *Machine) UpdateVMQuotas(id string) error {
	quota, err := carton.NewQuota(m.AccountId, id)
	if err != nil {
		return err
	}

	if m.Status == constants.StatusLaunching || m.Status == constants.StatusRunning {
			quota.Status = "activated"
			quota.AllocatedTo = m.CartonId
	} else if m.Status == constants.StatusDestroying {
		  quota.Status = "deactivated"
			quota.AllocatedTo = ""
	}

	return quota.Update()
}
