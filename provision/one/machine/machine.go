package machine

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	nsqp "github.com/crackcomm/nsqueue/producer"
	"github.com/megamsys/opennebula-go/compute"
	"github.com/megamsys/vertice/carton"
	"github.com/megamsys/vertice/events"
	"github.com/megamsys/vertice/events/alerts"
	"github.com/megamsys/vertice/meta"
	"github.com/megamsys/vertice/provision"
	"github.com/megamsys/vertice/provision/one/cluster"
)

type OneProvisioner interface {
	Cluster() *cluster.Cluster
}

type Machine struct {
	Name       string
	Id         string
	CartonId   string
	AccountsId string
	Level      provision.BoxLevel
	SSH        provision.BoxSSH
	Image      string
	Routable   bool
	Status     provision.Status
}

type CreateArgs struct {
	Commands    []string
	Box         *provision.Box
	Compute     provision.BoxCompute
	Deploy      bool
	Provisioner OneProvisioner
}

func (m *Machine) Create(args *CreateArgs) error {
	log.Infof("  creating machine in one (%s, %s)", m.Name, m.Image)

	opts := compute.VirtualMachine{
		Name:   m.Name,
		Image:  m.Image,
		Cpu:    string(args.Box.GetCpushare()), //ugly, compute has the info.
		Memory: string(args.Box.GetMemory()),
		HDD:    string(args.Box.GetHDD()),
		ContextMap: map[string]string{compute.ASSEMBLY_ID: args.Box.CartonId,
			compute.ASSEMBLIES_ID: args.Box.CartonsId},
	}

	//m.addEnvsToContext(m.BoxEnvs, &vm)

	_, _, err := args.Provisioner.Cluster().CreateVM(opts)
	if err != nil {
		return err
	}
	return nil
}

func (m *Machine) Remove(p OneProvisioner) error {
	log.Debugf("  removing machine in one (%s)", m.Name)
	opts := compute.VirtualMachine{
		Name: m.Name,
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
	mi[alerts.VERTNAME] = m.Name
	mi[alerts.COST] = "0.1"
	newEvent := events.NewMulti(
		[]*events.Event{
			&events.Event{
				AccountsId:  m.AccountsId,
				EventAction: alerts.DEDUCT,
				EventType:   events.EventBill,
				EventData:   events.EventData{M: mi},
				Timestamp:   time.Now().Local(),
			},
		})
	return newEvent.Write()
}

func (m *Machine) LifecycleOps(p OneProvisioner, action string) error {
	log.Debugf("  %s machine in one (%s)", action, m.Name)
	opts := compute.VirtualMachine{
		Name: m.Name,
	}
	err := p.Cluster().VM(opts, action)
	if err != nil {
		return err
	}
	return nil
}

//it possible to have a Notifier interface that does this, duck typed b y Assembly, Components.
func (m *Machine) SetStatus(status provision.Status) error {
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

//just publish a message stateup to the machine.
func (m *Machine) ChangeState(status provision.Status) error {
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
	fmt.Fprintf(w, "\nlogs nirvana ! machine %s ", m.Name)
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
