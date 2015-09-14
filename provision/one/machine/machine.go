// Copyright 2015 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package machine

import (
	"crypto"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/megamsys/megamd/db"
	"github.com/megamsys/megamd/provision"
	"github.com/megamsys/opennebula-go"
)

type Machine struct {
	Name          string
	Compute       BoxCompute
	ComponentId   string
	BoxEnvs       string
	Image         string
	Envs          string
	BuildingImage provision.Status
}

func (m *Machine) Available() bool {
	return c.Status == provision.StatusStarted.String() ||
		c.Status == provision.StatusStarting.String()
}

type CreateArgs struct {
	ImageID          string
	Commands         []string
	Box              provision.Box
	Deploy           bool
	Provisioner      OneProvisioner
	DestinationHosts []string
}

func (m *Machine) Create(args *CreateArgs) error {
	log.Debugf("creating machine %s in one with %s", m.Name, *args.ImageID)

	vm := compute.VirtualMachine{
		Name:         m.Name,
		TemplateName: m.Image,
		Cpu:          m.Compute.Cpushare,
		Memory:       m.Compute.Memory,
		Assembly_id:  m.box.AssemblyId,
		Client:       &p.client,
	}

	m.addEnvsToConfig(m.BoxEnvs, &vm)

	_, err := vm.Create()

	if err != nil {
		log.Errorf("error on creating machine in one %s - %s", m.Name, err)
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

func (m *Machine) SetStatus(status provision.Status) error {
	log.Debugf("setting status of machine %s to %s", m.Name, status.String())

	comp := carton.NewComponent(m.ComponentId)
	comp.SetStatus(status.String())

	if err != nil {
		log.Errorf("error on updating machine into riak %s - %s", m.ComponentId, err)
		return nil, err
	}

}

func (m *Machine) Remove(p OneProvisioner) error {
	log.Debugf("removing machine %s in one", m.Name)

	vm := compute.VirtualMachine{
		Name:   m.Name,
		Client: &p.client,
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
