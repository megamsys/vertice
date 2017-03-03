package cluster

import (
	"encoding/xml"
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/libgo/safe"
	constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/opennebula-go/api"
	"github.com/megamsys/opennebula-go/compute"
	"github.com/megamsys/opennebula-go/disk"
	"github.com/megamsys/opennebula-go/images"
	"github.com/megamsys/opennebula-go/snapshot"
	"github.com/megamsys/opennebula-go/template"
	"github.com/megamsys/opennebula-go/virtualmachine"
	"net"
	"net/url"
	"strconv"
	"time"
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

func (c *Cluster) GetVM(opts virtualmachine.Vnc, region string) (*virtualmachine.VM, error) {

	node, err := c.getNodeRegion(region)
	if err != nil {
		return nil, err
	}

	opts.T = node.Client
	res, err := opts.GetVm()
	if err != nil {
		return nil, wrapErrorWithCmd(node, err, "GetVM")
	}

	return res, err
}

// DestroyVM kills a vm, returning an error in case of failure.
func (c *Cluster) DestroyVM(opts compute.VirtualMachine) error {

	node, err := c.getNodeRegion(opts.Region)
	if err != nil {
		return err
	}

	opts.T = node.Client

	_, err = opts.Delete()
	if err != nil {
		return wrapErrorWithCmd(node, err, "DestroyVM")
	}

	return nil
}

// DestroyVM kills a vm, returning an error in case of failure.
func (c *Cluster) ForceDestoryVM(opts compute.VirtualMachine) error {

	node, err := c.getNodeRegion(opts.Region)
	if err != nil {
		return err
	}

	opts.T = node.Client

	_, err = opts.RecoverDelete()
	if err != nil {
		return wrapErrorWithCmd(node, err, "DestroyVM")
	}

	return nil
}

func (c *Cluster) VM(opts compute.VirtualMachine, action string) error {
	switch action {
	case constants.START:
		return c.StartVM(opts)
	case constants.STOP:
		return c.StopVM(opts)
	case constants.RESTART:
		return c.RestartVM(opts)
	case constants.HARD_STOP:
		return c.ForceStopVM(opts)
	case constants.HARD_RESTART:
		return c.ForceRestartVM(opts)
	default:
		return nil
	}
}
func (c *Cluster) StartVM(opts compute.VirtualMachine) error {

	node, err := c.getNodeRegion(opts.Region)
	if err != nil {
		return err
	}

	opts.T = node.Client

	_, err = opts.Resume()
	if err != nil {
		return wrapErrorWithCmd(node, err, "StartVM")
	}

	return nil
}

func (c *Cluster) RestartVM(opts compute.VirtualMachine) error {

	node, err := c.getNodeRegion(opts.Region)
	if err != nil {
		return err
	}

	opts.T = node.Client

	_, err = opts.Reboot()
	if err != nil {
		return wrapErrorWithCmd(node, err, "RebootVM")
	}

	return nil
}

func (c *Cluster) StopVM(opts compute.VirtualMachine) error {

	node, err := c.getNodeRegion(opts.Region)
	if err != nil {
		return err
	}

	opts.T = node.Client

	_, err = opts.PoweroffHard()
	if err != nil {
		return wrapErrorWithCmd(node, err, "StopVM")
	}

	return nil
}

func (c *Cluster) ForceRestartVM(opts compute.VirtualMachine) error {

	node, err := c.getNodeRegion(opts.Region)
	if err != nil {
		return err
	}

	opts.T = node.Client

	_, err = opts.RebootHard()
	if err != nil {
		return wrapErrorWithCmd(node, err, "RebootVM")
	}

	return nil
}

func (c *Cluster) ForceStopVM(opts compute.VirtualMachine) error {

	node, err := c.getNodeRegion(opts.Region)
	if err != nil {
		return err
	}

	opts.T = node.Client

	_, err = opts.Poweroff()
	if err != nil {
		return wrapErrorWithCmd(node, err, "StopVM")
	}

	return nil
}

func (c *Cluster) getNodeRegion(region string) (node, error) {
	return c.getNode(func(s Storage) (Node, error) {
		return s.RetrieveNode(region)
	})
}

func (c *Cluster) SaveDiskImage(opts compute.Image) (string, error) {
	node, err := c.getNodeRegion(opts.Region)
	if err != nil {
		return "", err
	}
	opts.T = node.Client

	res, err := opts.DiskSaveAs()
	if err != nil {
		return "", wrapErrorWithCmd(node, err, "CreateImage")
	}
	imageId := res.(string)
	return imageId, nil
}

func (c *Cluster) RemoveImage(opts compute.Image) error {
	node, err := c.getNodeRegion(opts.Region)
	if err != nil {
		return err
	}
	opts.T = node.Client

	_, err = opts.RemoveImage()
	if err != nil {
		return wrapErrorWithCmd(node, err, "DeleteSnap")
	}

	return nil
}

func (c *Cluster) IsImageReady(v *images.Image, region string) error {
	node, err := c.getNodeRegion(region)
	if err != nil {
		return err
	}
	v.T = node.Client
	err = safe.WaitCondition(30*time.Minute, 20*time.Second, func() (bool, error) {
		res, err := v.Show()
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

func (c *Cluster) SnapVMDisk(opts snapshot.Snapshot, region string) (string, error) {
	node, err := c.getNodeRegion(region)
	if err != nil {
		return "", err
	}
	opts.T = node.Client

	res, err := opts.CreateSnapshot()
	if err != nil {
		return "", wrapErrorWithCmd(node, err, "CreateSnap")
	}
	imageId := res.(string)
	return imageId, nil
}

func (c *Cluster) RemoveSnap(opts snapshot.Snapshot, region string) error {
	node, err := c.getNodeRegion(region)
	if err != nil {
		return err
	}
	opts.T = node.Client

	_, err = opts.DeleteSnapshot()
	if err != nil {
		return wrapErrorWithCmd(node, err, "DeleteSnap")
	}

	return nil
}

func (c *Cluster) RestoreSnap(opts snapshot.Snapshot, region string) error {
	node, err := c.getNodeRegion(region)
	if err != nil {
		return err
	}
	opts.T = node.Client

	_, err = opts.RevertSnapshot()
	if err != nil {
		return wrapErrorWithCmd(node, err, "RestoreSnap")
	}

	return nil
}

func (c *Cluster) IsSnapReady(v *images.Image, region string) error {
	node, err := c.getNodeRegion(region)
	if err != nil {
		return err
	}
	v.T = node.Client
	err = safe.WaitCondition(10*time.Minute, 10*time.Second, func() (bool, error) {
		res, err := v.Show()
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

func (c *Cluster) GetDiskId(vd *disk.VmDisk, region string) ([]int, error) {
	var a []int
	node, err := c.getNodeRegion(region)
	if err != nil {
		return a, err
	}
	vd.T = node.Client

	dsk, err := vd.ListDisk()
	if err != nil {
		return a, err
	}

	a = dsk.GetDiskIds()
	return a, nil
}

func (c *Cluster) AttachDisk(v *disk.VmDisk, region string) error {

	node, err := c.getNodeRegion(region)
	if err != nil {
		return err
	}

	v.T = node.Client

	_, err = v.AttachDisk()
	if err != nil {
		return wrapErrorWithCmd(node, err, "AttachDisk")
	}
	return nil
}

func (c *Cluster) DetachDisk(v *disk.VmDisk, region string) error {

	node, err := c.getNodeRegion(region)
	if err != nil {
		return err
	}

	v.T = node.Client

	_, err = v.DetachDisk()
	if err != nil {
		return wrapErrorWithCmd(node, err, "DetachDisk")
	}

	return nil
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

func (c *Cluster) InstantiateVM(opts *template.UserTemplate, vname, throttle, region string) (string, error) {
	var (
		addr string
		vmid string
		err  error
	)
	maxTries := 5
	nodlist, err := c.Nodes()
	for _, v := range nodlist {
		if v.Metadata[api.ONEZONE] == region {
			addr = v.Address
			if v.Metadata[api.VCPU_PERCENTAGE] != "" {
				opts.Template.Cpu = cpuThrottle(v.Metadata[api.VCPU_PERCENTAGE], opts.Template.Cpu)
			} else {
				opts.Template.Cpu = cpuThrottle(throttle, opts.Template.Cpu)
			}
		}
	}

	if addr == "" {
		return vmid, fmt.Errorf("%s", cmd.Colorfy("Unavailabldd region ( "+region+" ) nodes (hint: start or beat it).\n", "red", "", ""))
	}
	userTemps := make([]*template.UserTemplate, 0)
	for ; maxTries > 0; maxTries-- {
		finalXML := template.UserTemplates{}
		finalXML.UserTemplate = append(userTemps, opts)
		finalData, err := xml.Marshal(finalXML)
		if err == nil {
			tmp := &template.TemplateReqs{
				TemplateName: opts.Template.Name,
				TemplateId:   opts.Id,
				TemplateData: string(finalData),
			}
			vmid, err = c.instantiateVMInNode(tmp, vname, region)
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
		return vmid, err
	}
	if err != nil {
		return vmid, fmt.Errorf("CreateVM: maximum number of tries exceeded, last error: %s", err.Error())
	}
	return vmid, err
}

func (c *Cluster) instantiateVMInNode(v *template.TemplateReqs, vmname, region string) (string, error) {

	node, err := c.getNodeRegion(region)
	if err != nil {
		return "", err
	}
	v.T = node.Client

	res, err := v.Instantiate(vmname)
	if err != nil {
		return "", wrapErrorWithCmd(node, err, "InstantiateVM")
	}

	return res.(string), nil
}

func (c *Cluster) ImagePersistent(opts images.Image, region string) error {
	node, err := c.getNodeRegion(region)
	if err != nil {
		return err
	}

	opts.T = node.Client
	_, err = opts.ChPersistent(false)
	return err
}

func (c *Cluster) GetImage(opts images.Image, region string) (*images.Image, error) {
	node, err := c.getNodeRegion(region)
	if err != nil {
		return nil, err
	}
	opts.T = node.Client
	return opts.Show()
}

func (c *Cluster) GetTemplate(region string) (*template.UserTemplate, error) {
	node, err := c.getNodeRegion(region)
	if err != nil {
		return nil, err
	}
	templateObj := &template.TemplateReqs{TemplateName: node.template, T: node.Client}
	res, err := templateObj.Get()
	if err != nil {
		return nil, err
	}
	return res[0], nil
}

func cpuThrottle(vcpu, cpu string) string {
	ThrottleFactor, _ := strconv.Atoi(vcpu)
	ICpu, _ := strconv.Atoi(cpu)
	realCPU := float64(ICpu) / float64(ThrottleFactor)
	//ugly, compute has the info.
	return strconv.FormatFloat(realCPU, 'f', 6, 64)
}

func (c *Cluster) AttachNics(rules map[string]string, vmid, region, storage string) error {
	vnets, err := c.getNics(rules, region, storage)
	if err != nil {
		return err
	}
	return c.attachNics(vnets, vmid, region)
}

func (c *Cluster) attachNics(vnets []string, vmid, region string) error {
	var failures []error
	node, err := c.getNodeRegion(region)
	if err != nil {
		return err
	}

	opts := virtualmachine.Vnc{VmId: vmid}
	opts.T = node.Client
	for _, vnet := range vnets {
		err := opts.AttachNic(vnet)
		err = safe.WaitCondition(1*time.Minute, 5*time.Second, func() (bool, error) {
			res, err := opts.GetVm()
			if err != nil {
				return false, err
			}
			return (res.State == int(virtualmachine.ACTIVE) && res.LcmState == int(virtualmachine.RUNNING)), nil
		})
		if err != nil {
			failures = append(failures, err)
			log.Debugf("  failed to attach nic (%s)", err)
		}
	}
	return nil
}

func (c *Cluster) DetachNics(net_ids []string, vmid, region string) error {
	var failures []error
	node, err := c.getNodeRegion(region)
	if err != nil {
		return err
	}

	opts := virtualmachine.Vnc{VmId: vmid}
	opts.T = node.Client
	for _, nid := range net_ids {
		id, _ := strconv.Atoi(nid)
		err := opts.DetachNic(id)
		err = safe.WaitCondition(1*time.Minute, 5*time.Second, func() (bool, error) {
			res, err := opts.GetVm()
			if err != nil {
				return false, err
			}
			return (res.State == int(virtualmachine.ACTIVE) && res.LcmState == int(virtualmachine.RUNNING)), nil
		})
		if err != nil {
			failures = append(failures, err)
			log.Debugf("  failed to attach nic (%s)", err)
		}
	}
	return nil
}

func (c *Cluster) getNics(rules map[string]string, region, storage string) ([]string, error) {
	vnets := make([]string, 0)
	var nics, nics_count map[string]string
	nodlist, err := c.Nodes()
	if err != nil {
		return vnets, err
	}
	for _, v := range nodlist {
		if v.Metadata[api.ONEZONE] == region {
			nics, nics_count = c.netAttachPolicy(v, rules, storage)
		}
	}

	for key, vnet := range nics {
		count, _ := strconv.Atoi(nics_count[key])
		for i := 0; i < count; i++ {
			vnets = append(vnets, vnet)
		}
	}
	if len(nics) == 0 {
		return vnets, fmt.Errorf("no networks available in this region (%s)", region)
	}

	return vnets, nil
}

// "rules":  [{"key":"ipv4private","value":"4"},{"key":"ipv6private","value":"2"}]
// if rules has key then store one vnet in nic[key] and  count of vnet in nic_count[key]
// nodeo.Clusters[Region.ClusterId] has values like map of {"ipv4private":"private-net4","ipv6private":"private-net6"}

func (c *Cluster) netAttachPolicy(nodeo Node, rules map[string]string, st string) (map[string]string, map[string]string) {
	var nic, nic_count map[string]string
	for id, cluster := range nodeo.Clusters {
		if cluster[constants.STORAGE_TYPE] == st && cluster[constants.VONE_CLOUD] != constants.TRUE {
			for nic_key, nic_value := range nodeo.Clusters[id] {
				if count, ok := rules[nic_key]; ok {
					nic[nic_key] = nic_value
					nic_count[nic_key] = count
				}
			}
			if len(nic) > 0 {
				return nic, nic_count
			}
		}
	}
	return nic, nic_count
}
