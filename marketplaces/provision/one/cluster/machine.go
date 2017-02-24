package cluster

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/libgo/safe"
	constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/opennebula-go/api"
	"github.com/megamsys/opennebula-go/compute"
	"github.com/megamsys/opennebula-go/images"
	"net"
	"net/url"
	"strconv"
	"time"
)

// CreateVM creates a vm in the specified node.
// It returns the vm, or an error, in case of failures.
const (
	START   = "start"
	STOP    = "stop"
	RESTART = "restart"
)

var ErrConnRefused = errors.New("connection refused")

func (c *Cluster) CreateVM(opts compute.VirtualMachine, throttle, storage string) (string, string, string, error) {

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
				opts.Vnets, opts.ClusterId = c.getVnets(v, opts.Vnets, storage)
				if v.Metadata[api.VCPU_PERCENTAGE] != "" {
					opts.Cpu = cpuThrottle(v.Metadata[api.VCPU_PERCENTAGE], opts.Cpu)
				} else {
					opts.Cpu = cpuThrottle(throttle, opts.Cpu)
				}
			}
		}

		if addr == "" {
			return addr, machine, vmid, fmt.Errorf("%s", cmd.Colorfy("Unavailable region ( "+opts.Region+" ) nodes (hint: start or beat it).\n", "red", "", ""))
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
	node, err := c.getNodeRegion(opts.Region)
	if err != nil {
		return "", "", err
	}

	if opts.ClusterId != "" {
		opts.TemplateName = node.template
	} else {
		opts.TemplateName = opts.Image
	}

	opts.T = node.Client

	res, err := opts.Create()
	if err != nil {
		return "", "", err
	}
	vmid := res.(string)
	return opts.Name, vmid, nil
}

func (c *Cluster) ImageCreate(opts images.Image, region string) (interface{}, error) {
	var ds string
	nodlist, err := c.Nodes()
	for _, v := range nodlist {
		if v.Metadata[api.ONEZONE] == region {
			ds = v.Metadata[constants.DATASTORE]
			if ds == "" {
				return ds, fmt.Errorf("%s", cmd.Colorfy("Datastore id is empty (hint: start or beat it).\n", "red", "", ""))
			}
			break
		}
	}

	if ds == "" {
		return ds, fmt.Errorf("%s", cmd.Colorfy("Unavailable region ( "+region+" ) nodes (hint: start or beat it).\n", "red", "", ""))
	}

	node, err := c.getNodeRegion(region)
	if err != nil {
		return nil, err
	}

	ds_id, err := strconv.Atoi(ds)
	if err != nil {
		return nil, wrapErrorWithCmd(node, err, "createimage")
	}
	opts.T = node.Client
	opts.DatastoreID = ds_id
	return opts.Create()
}

func (c *Cluster) IsImageReady(v *images.Image, region string) error {
	node, err := c.getNodeRegion(region)
	if err != nil {
		return err
	}
	v.T = node.Client
	err = safe.WaitCondition(10*time.Minute, 10*time.Second, func() (bool, error) {
		res, err := v.ImageShow()
		if err != nil || res.State_string() == "failure" {
			return false, fmt.Errorf("fails to create snapshot")
		}
		return (res.State_string() == "ready"), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Cluster) getNodeRegion(region string) (node, error) {
	return c.getNode(func(s Storage) (Node, error) {
		return s.RetrieveNode(region)
	})
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
