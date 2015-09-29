package machine

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/amqp"
	"github.com/megamsys/megamd/carton"
	"github.com/megamsys/megamd/meta"
	"github.com/megamsys/megamd/provision"
	"github.com/megamsys/megamd/provision/one/cluster"
	"github.com/megamsys/opennebula-go/compute"
)

type OneProvisioner interface {
	Cluster() *cluster.Cluster
}

type Machine struct {
	Name     string
	Id       string
	CartonId string
	Level    provision.BoxLevel
	Image    string
	Routable bool
}

type CreateArgs struct {
	Commands    []string
	Box         *provision.Box
	Compute     provision.BoxCompute
	Deploy      bool
	Provisioner OneProvisioner
}

func (m *Machine) Create(args *CreateArgs) error {
	log.Debugf("creating machine %s in one with %s", m.Name, m.Image)

	opts := compute.VirtualMachine{
		Name:   m.Name,
		Image:  m.Image,
		Cpu:    args.Compute.Cpushare,
		Memory: args.Compute.Memory,
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
	log.Debugf("removing machine %s in one", m.Name)
	opts := compute.VirtualMachine{
		Name: m.Name,
	}

	err := p.Cluster().DestroyVM(opts)
	if err != nil {
		return err
	}
	return nil
}

//it possible to have a Notifier interface that does this, duck typed by Assembly, Components.
func (m *Machine) SetStatus(status provision.Status) error {
	log.Debugf("setting status of machine %s %s to %s", m.Id, m.Name, status.String())

	if asm, err := carton.NewAssembly(m.CartonId); err != nil {
		return err
	} else if err = asm.SetStatus(status); err != nil {
		return err
	}

	if m.Level == provision.BoxSome {
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
	log.Debugf("change state of machine %s %s to %s", m.Id, m.Name, status.String())

	p, err := amqp.NewRabbitMQ(meta.MC.AMQP, m.Name)
	if err != nil {
		return err
	}

	jsonMsg, err := json.Marshal(
		carton.Requests{
			Action:    status.String(),
			Category:  carton.STATE,
			CreatedAt: time.Now().String(),
		})

	if err != nil {
		return err
	}

	if err := p.Pub(jsonMsg); err != nil {
		return err
	}
	return nil

}

//if there is a file or something to be created, do it here.
func (m *Machine) Logs(p OneProvisioner, w io.Writer) error {
	log.Debugf("hook machine %s logs", m.Name)
	fmt.Fprintf(w, "\nhook machine %s logs", m.Name)
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
