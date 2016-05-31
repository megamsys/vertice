package cluster

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"encoding/json"
	 "strings"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/opennebula-go/compute"
	"github.com/megamsys/opennebula-go/virtualmachine"
)

// CreateVM creates a vm in the specified node.
// It returns the vm, or an error, in case of failures.
const (
	START   = "start"
	STOP    = "stop"
	RESTART = "restart"
)


var ErrConnRefused = errors.New("connection refused")

func (c *Cluster) CreateVM(opts compute.VirtualMachine) (string, string, string, error) {
	var (
		addr    string
		machine string
		vmid    string
		err     error
	)
	maxTries := 5
	for ; maxTries > 0; maxTries-- {

		nodlist, err := c.Nodes()

		if err != nil || len(nodlist) <= 0 {
			return addr, machine, vmid, fmt.Errorf("%s", cmd.Colorfy("Unavailable nodes (hint: start or beat it).\n", "red", "", ""))
		} else {
			addr = nodlist[0].Address
		}

		if err == nil {
			machine, vmid, err = c.createVMInNode(opts, addr)
			if err == nil {
				c.handleNodeSuccess(addr)
				break
			}
			log.Errorf("  > Trying... %s", addr)
		}
		shouldIncrementFailures := false
		isCreateMachineErr := false
		baseErr := err
		if nodeErr, ok := baseErr.(OneNodeError); ok {
			isCreateMachineErr = nodeErr.cmd == "createVM"
			baseErr = nodeErr.BaseError()
		}
		if urlErr, ok := baseErr.(*url.Error); ok {
			baseErr = urlErr.Err
		}
		_, isNetErr := baseErr.(*net.OpError)
		if isNetErr || isCreateMachineErr || baseErr == ErrConnRefused {
			shouldIncrementFailures = true
		}
		c.handleNodeError(addr, err, shouldIncrementFailures)
		return addr, machine, vmid, err
	}
	if err != nil {
		return addr, machine, vmid,fmt.Errorf("CreateVM: maximum number of tries exceeded, last error: %s", err.Error())
	}
	return addr, machine, vmid, err
}

//create a vm in a node.
func (c *Cluster) createVMInNode(opts compute.VirtualMachine, nodeAddress string) ( string, string, error) {
	 node, err := c.getNodeByAddr(nodeAddress)
	if err != nil {
		return "", "",err
	}
	opts.TemplateName = node.template
	opts.T = node.Client

	res, err := opts.Create()
	 b, err :=json.Marshal(res)
	 if err != nil {
		 return "", "", err
	 }
	  str :=string(b)
		 spstr :=strings.Split(str,",")
		 vmres := spstr[1]
	if err != nil {
		return "", "", wrapErrorWithCmd(node, err, "createVM")
	}

	return opts.Name, vmres, nil
}



func (c *Cluster) GetIpPort(opts virtualmachine.Vnc) ( string, string, error) {

	var (
		addr string
		hostip string
		vncport string
	)
	nodlist, err := c.Nodes()

	if err != nil || len(nodlist) <= 0 {
		return hostip, vncport, fmt.Errorf("%s", cmd.Colorfy("Unavailable nodes (hint: start or beat it).\n", "red", "", ""))
	} else {
		addr = nodlist[0].Address
	}

	node, err := c.getNodeByAddr(addr)
	if err != nil {
		return "", "",err
	}
	//opts.TemplateName = node.template
	opts.T = node.Client

	res, err := opts.GetVm()
	 vncport = res.GetPort()
	 hostip = res.GetHostIp()
	if err != nil {
		return "", "", wrapErrorWithCmd(node, err, "createVM")
	}

	return hostip, vncport, nil
}
// DestroyVM kills a vm, returning an error in case of failure.
func (c *Cluster) DestroyVM(opts compute.VirtualMachine) error {
	var (
		addr string
	)
	nodlist, err := c.Nodes()

	if err != nil || len(nodlist) <= 0 {
		return fmt.Errorf("%s", cmd.Colorfy("Unavailable nodes (hint: start or beat it).\n", "red", "", ""))
	} else {
		addr = nodlist[0].Address
	}

	node, err := c.getNodeByAddr(addr)
	if err != nil {
		return err
	}
	opts.T = node.Client

	_, err = opts.Delete()
	if err != nil {
		return wrapError(node, err)
	}
	return nil
}

func (c *Cluster) VM(opts compute.VirtualMachine, action string) error {
	switch action {
	case START:
		return c.StartVM(opts)
	case STOP:
		return c.StopVM(opts)
	case RESTART:
		return c.RestartVM(opts)
	default:
		return nil
	}
}
func (c *Cluster) StartVM(opts compute.VirtualMachine) error {
	var (
		addr string
	)
	nodlist, err := c.Nodes()

	if err != nil || len(nodlist) <= 0 {
		return fmt.Errorf("%s", cmd.Colorfy("Unavailable nodes (hint: start or beat it).\n", "red", "", ""))
	} else {
		addr = nodlist[0].Address
	}

	node, err := c.getNodeByAddr(addr)
	if err != nil {
		return err
	}
	opts.T = node.Client

	_, err = opts.Resume()
	if err != nil {
		return wrapError(node, err)
	}
	return nil
}

func (c *Cluster) RestartVM(opts compute.VirtualMachine) error {
	var (
		addr string
	)
	nodlist, err := c.Nodes()

	if err != nil || len(nodlist) <= 0 {
		return fmt.Errorf("%s", cmd.Colorfy("Unavailable nodes (hint: start or beat it).\n", "red", "", ""))
	} else {
		addr = nodlist[0].Address
	}

	node, err := c.getNodeByAddr(addr)
	if err != nil {
		return err
	}
	opts.T = node.Client

	_, err = opts.Reboot()
	if err != nil {
		return wrapError(node, err)
	}
	return nil
}

func (c *Cluster) StopVM(opts compute.VirtualMachine) error {
	var (
		addr string
	)
	nodlist, err := c.Nodes()
	if err != nil || len(nodlist) <= 0 {
		return fmt.Errorf("%s", cmd.Colorfy("Unavailable nodes (hint: start or beat it).\n", "red", "", ""))
	} else {
		addr = nodlist[0].Address
	}

	node, err := c.getNodeByAddr(addr)
	if err != nil {
		return err
	}
	opts.T = node.Client

	_, err = opts.Poweroff()
	if err != nil {
		return wrapError(node, err)
	}
	return nil
}

func (c *Cluster) getNodeByAddr(addr string) (node, error) {
	return c.getNode(func(s Storage) (Node, error) {
		return s.RetrieveNode(addr)
	})
}
