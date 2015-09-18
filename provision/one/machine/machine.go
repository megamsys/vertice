// Copyright 2015 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package machine

import (
	"fmt"
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/megamd/carton"
	"github.com/megamsys/megamd/provision"
	"github.com/megamsys/opennebula-go/api"
	"github.com/megamsys/opennebula-go/compute"
)

type OneProvisioner interface {
	Cluster() *api.Rpc
}

type Machine struct {
	Name     string
	Id       string
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

	vm := compute.VirtualMachine{
		Name:         m.Name,
		TemplateName: m.Image,
		Cpu:          args.Compute.Cpushare,
		Memory:       args.Compute.Memory,
		Assembly_id:  args.Box.CartonId,
		Client:       args.Provisioner.Cluster(),
	}

	//m.addEnvsToContext(m.BoxEnvs, &vm)

	_, err := vm.Create()

	if err != nil {
		return err
	}
	//do you want to update something.
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

//it possible to have a Notifier interface that does this, duck typed by Assembly, Components.
func (m *Machine) SetStatus(status provision.Status) error {
	log.Debugf("setting status of machine %s %s to %s", m.Id, m.Name, status.String())

	switch m.Level {
	case provision.BoxSome:
		if comp, err := carton.NewComponent(m.Id); err != nil {
			return err
		} else if err = comp.SetStatus(status); err != nil {
			return nil
		}
		return nil
	case provision.BoxNone:
		if asm, err := carton.NewAssembly(m.Id); err != nil {
			return err
		} else if err = asm.SetStatus(status); err != nil {
			return nil
		}
		return nil
	default:
		return nil
	}

	return nil
}

func (m *Machine) Remove(p OneProvisioner) error {
	log.Debugf("removing machine %s in one", m.Name)

	vm := compute.VirtualMachine{
		Name:   m.Name,
		Client: p.Cluster(),
	}

	_, err := vm.Delete()

	if err != nil {
		log.Errorf("error on deleting machine in one %s - %s", m.Name, err)
		return err
	}

	//do you want to update something here ?
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

func (m *Machine) Logs(p OneProvisioner, w io.Writer) error {
	log.Debugf("hook machine %s logs", m.Name)
	fmt.Fprintf(w, "\nhook machine %s logs", m.Name)
	//if there is a file or something to be created, do it here.
	return nil
}
