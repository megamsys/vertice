package cluster

import (
	"errors"
	"fmt"
	"net"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/opennebula-go/api"
	"github.com/megamsys/opennebula-go/compute"
)

// CreateVM creates a vm in the specified node.
// It returns the vm, or an error, in case of failures.
func (c *Cluster) CreateVM(opts compute.VirtualMachine) (string, string, error) {
	var (
		addr    string
		machine string
		err     error
	)
	maxTries := 5
	for ; maxTries > 0; maxTries-- {

		nodlist, err := c.Nodes()

		if err != nil || len(nodlist) <= 0 {
			return addr, machine, fmt.Errorf("%s\n%s", cmd.Colorfy("Nodes are not available to launch machines.\n%s", "red", "", ""), err)
		} else {
			addr = nodlist[0].Address
		}

		if err == nil {
			machine, err = c.createVMInNode(opts, addr)
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
		if isNetErr || isCreateMachineErr || baseErr == api.ErrConnRefused {
			shouldIncrementFailures = true
		}
		c.handleNodeError(addr, err, shouldIncrementFailures)
		return addr, machine, err
	}
	if err != nil {
		return addr, machine, fmt.Errorf("CreateVM: maximum number of tries exceeded, last error: %s", err.Error())
	}
	return addr, machine, err
}

//create a vm in a node.
func (c *Cluster) createVMInNode(opts compute.VirtualMachine, nodeAddress string) (string, error) {
	node, err := c.getNodeByAddr(nodeAddress)
	if err != nil {
		return "", err
	}
	opts.TemplateName = node.template
	opts.Client = node.Client

	_, err = opts.Create()

	if err != nil {
		return "", wrapErrorWithCmd(node, err, "createVM")
	}

	return opts.Name, nil
}

// DestroyVM kills a vm, returning an error in case of failure.
func (c *Cluster) DestroyVM(opts compute.VirtualMachine) error {
	var (
		addr string
	)
	if nodlist, err := c.Nodes(); err != nil {
		return errors.New("DeleteVM needs a non empty node addr")
	} else {
		addr = nodlist[0].Address
	}

	node, err := c.getNodeByAddr(addr)
	if err != nil {
		return err
	}
	opts.Client = node.Client

	_, err = opts.Delete()
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
