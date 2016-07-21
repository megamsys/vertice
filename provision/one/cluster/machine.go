package cluster

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/opennebula-go/api"
	"github.com/megamsys/opennebula-go/compute"
	"github.com/megamsys/opennebula-go/virtualmachine"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// CreateVM creates a vm in the specified node.
// It returns the vm, or an error, in case of failures.
const (
	START   = "start"
	STOP    = "stop"
	RESTART = "restart"
)

var ErrConnRefused = errors.New("connection refused")

func (c *Cluster) CreateVM(opts compute.VirtualMachine, t string) (string, string, string, error) {
	var (
		addr    string
		machine string
		vmid    string
		err     error
	)
	maxTries := 5
	for ; maxTries > 0; maxTries-- {

		nodlist, err := c.Nodes()

		for _, v := range nodlist {
			if v.Metadata[api.ONEZONE] == opts.Region {
				addr = v.Address
				opts.Vnets, opts.ClusterId = c.getVnets(v, opts.Vnets)
				if v.Metadata[api.VCPU_PERCENTAGE] != "" {
					opts.Cpu = cpuThrottle(v.Metadata[api.VCPU_PERCENTAGE], opts.Cpu)
				} else {
					opts.Cpu = cpuThrottle(t, opts.Cpu)
				}
			}
		}

		if addr == "" {
			return addr, machine, vmid, fmt.Errorf("%s", cmd.Colorfy("Unavailable nodes (hint: start or beat it).\n", "red", "", ""))
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
		return addr, machine, vmid, fmt.Errorf("CreateVM: maximum number of tries exceeded, last error: %s", err.Error())
	}
	return addr, machine, vmid, err
}

//create a vm in a node.
func (c *Cluster) createVMInNode(opts compute.VirtualMachine, nodeAddress string) (string, string, error) {
	node, err := c.getNodeByAddr(nodeAddress)
	if err != nil {
		return "", "", err
	}
	opts.TemplateName = node.template
	opts.T = node.Client

	res, err := opts.Create()
	b, err := json.Marshal(res)

	if err != nil {
		return "", "", err
	}
	str := string(b)
	spstr := strings.Split(str, ",")
	vmres := spstr[1]
	if err != nil {
		return "", "", wrapErrorWithCmd(node, err, "createVM")
	}

	return opts.Name, vmres, nil
}

func (c *Cluster) GetIpPort(opts virtualmachine.Vnc, region string) (string, string, error) {

	addr, err := c.getRegion(region)
	if err != nil {
		return "", "", err
	}

	node, err := c.getNodeByAddr(addr)
	if err != nil {
		return "", "", err
	}
	//opts.TemplateName = node.template
	opts.T = node.Client

	res, err := opts.GetVm()
	vnchost := res.GetHostIp()
	vncport := res.GetPort()

	if err != nil {
		return "", "", wrapErrorWithCmd(node, err, "createVM")
	}

	return vnchost, vncport, nil
}

// DestroyVM kills a vm, returning an error in case of failure.
func (c *Cluster) DestroyVM(opts compute.VirtualMachine) error {

	addr, err := c.getRegion(opts.Region)
	if err != nil {
		return err
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

	addr, err := c.getRegion(opts.Region)
	if err != nil {
		return err
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

	addr, err := c.getRegion(opts.Region)
	if err != nil {
		return err
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

	addr, err := c.getRegion(opts.Region)
	if err != nil {
		return err
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

func (c *Cluster) SnapVMDisk(opts compute.VirtualMachine) error {

	addr, err := c.getRegion(opts.Region)
	if err != nil {
		return err
	}

	node, err := c.getNodeByAddr(addr)
	if err != nil {
		return err
	}
	opts.T = node.Client

	_, err = opts.DiskSnap()
	if err != nil {
		return wrapError(node, err)
	}
	return nil
}

func cpuThrottle(vcpu, cpu string) string {
	ThrottleFactor, _ := strconv.Atoi(vcpu)
	cpuThrottleFactor := float64(ThrottleFactor)
	ICpu, _ := strconv.Atoi(cpu)
	throttle := float64(ICpu)
	realCPU := throttle / cpuThrottleFactor
	cpu = strconv.FormatFloat(realCPU, 'f', 6, 64) //ugly, compute has the info.
	return cpu
}

func (c *Cluster) getRegion(region string) (string, error) {
	var (
		addr string
	)
	nodlist, err := c.Nodes()
	if err != nil {
		addr = ""
	}
	for _, v := range nodlist {
		if v.Metadata[api.ONEZONE] == region {
			addr = v.Address
		}
	}

	if addr == "" {
		return addr, fmt.Errorf("%s", cmd.Colorfy("Unavailable nodes (hint: start or beat it).\n", "red", "", ""))
	}

	return addr, nil
}
