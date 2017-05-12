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
	"github.com/megamsys/opennebula-go/disk"
	"github.com/megamsys/opennebula-go/images"
	"github.com/megamsys/opennebula-go/snapshot"
	"github.com/megamsys/opennebula-go/template"
	"github.com/megamsys/opennebula-go/virtualmachine"
	"github.com/megamsys/vertice/carton"
	lb "github.com/megamsys/vertice/logbox"
	"github.com/megamsys/vertice/meta"
	"github.com/megamsys/vertice/provision"
	"github.com/megamsys/vertice/provision/one/cluster"

	mk "github.com/megamsys/vertice/marketplaces"
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
	PublicUrl    string
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

func (m *Machine) create(args *CreateArgs) (compute.VirtualMachine, *carton.Assembly, error) {
	asm, err := carton.NewAssembly(m.CartonId, m.AccountId, "")
	if err != nil {
		return compute.VirtualMachine{}, nil, err
	}
	flv, err := carton.GetFlavor(m.AccountId, asm.FlavorId())
	if err != nil {
		return compute.VirtualMachine{}, nil, err
	}
	opts := compute.VirtualMachine{
		Name:       m.Name,
		Image:      m.Image,
		Region:     args.Box.Region,
		Cpu:        strconv.FormatInt(int64(args.Box.GetCpushare()), 10),
		Memory:     strconv.FormatInt(int64(args.Box.GetMemory()), 10),
		HDD:        strconv.FormatInt(int64(args.Box.GetHDD()), 10),
		CpuCost:    flv.GetCpuCost(),
		MemoryCost: flv.GetMemoryCost(),
		HDDCost:    flv.GetHDDCost(),
		ContextMap: map[string]string{compute.ASSEMBLY_ID: args.Box.CartonId, compute.ORG_ID: args.Box.OrgId,
			compute.ASSEMBLIES_ID: args.Box.CartonsId, compute.ACCOUNTS_ID: args.Box.AccountId, compute.API_KEY: args.Box.ApiArgs.Api_Key, constants.QUOTA_ID: args.Box.QuotaId},
		Vnets: args.Box.Vnets,
	}
	opts.VCpu = opts.Cpu
	if strings.Contains(args.Box.Tosca, "freebsd") {
		opts.Files = "/detio/freebsd/init.sh"
	}
	return opts, asm, nil
}

func (m *Machine) Create(args *CreateArgs) error {
	nics := make([]*template.NIC, 0)
	opts, asm, err := m.create(args)
	_, _, vmid, err := args.Provisioner.Cluster().CreateVM(opts, m.VCPUThrottle, m.StorageType, nics)
	if err != nil {
		return err
	}
	m.VMId = vmid
	var id = make(map[string][]string)
	id[carton.INSTANCE_ID] = []string{m.VMId}
	if err = asm.NukeAndSetOutputs(id); err != nil {
		return err
	}
	return nil
}

func (m *Machine) CreateBackupVM(args *CreateArgs) error {
	var ips = make(map[string][]string)
	asm, err := carton.NewAssembly(m.CartonId, m.AccountId, "")
	if err != nil {
		return err
	}
	bk, err := carton.GetBackup(asm.Inputs.Match("backup_id"), m.AccountId)
	if err != nil {
		return err
	}

	nics := []string{constants.PUBLICIPV4, constants.PRIVATEIPV4, constants.PUBLICIPV6, constants.PRIVATEIPV6}
	for _, nic := range nics {
		if ip := bk.Outputs.Match(nic); ip != "" {
			t := strings.Split(ip, ",")
			if len(t) > 0 {
				ips[nic] = t
			}
		}
	}

	opts, asm, err := m.create(args)
	if err != nil {
		return err
	}

	res, err := args.Provisioner.Cluster().GetIpsNetwork(m.Region, ips)
	if err != nil {
		err = m.SetStatusErr(constants.StatusNetworkUnavail, err)
		opts.ForceNetwork = false
	} else {
		opts.ForceNetwork = true
	}

	_, _, vmid, err := args.Provisioner.Cluster().CreateVM(opts, m.VCPUThrottle, m.StorageType, res)
	if err != nil {
		return err
	}
	m.VMId = vmid
	var id = make(map[string][]string)
	id[carton.INSTANCE_ID] = []string{m.VMId}
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
		msg := fmt.Sprintf("quota payment pending for the machine (%s) ", b.GetFullName())
		carton.DoneNotify(b, w, alerts.INSUFFICIENT_FUND, msg)
		log.Debugf(msg)
		return fmt.Errorf("credit balance insufficient")
	}

	return nil
}

func (m *Machine) CheckQuotaState(b *provision.Box, w io.Writer) error {
	quota, err := carton.NewQuota(m.AccountId, b.QuotaId)
	if err != nil {
		return err
	}
	if strings.ToLower(quota.Status) != "paid" {
		msg := fmt.Sprintf("quota payment pending for the machine (%s) ", b.GetFullName())
		carton.DoneNotify(b, w, alerts.QUOTA_UNPAID, msg)
		log.Debugf(msg)
		return fmt.Errorf("quota state unpaid")
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

	err = safe.WaitCondition(30*time.Minute, 20*time.Second, func() (bool, error) {
		_ = asm.Trigger_event(utils.Status(constants.StatusWaitUntill))
		res, err = args.Provisioner.Cluster().GetVM(opts, m.Region)
		if err != nil {
			return false, err
		}
		if res.State == int(virtualmachine.DONE) {
			return false, fmt.Errorf("VM deleted while machine deploying")
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

func (m *Machine) WaitUntillVMState(p OneProvisioner, vm virtualmachine.VmState, lcm virtualmachine.LcmState) error {
	opts := virtualmachine.Vnc{VmId: m.VMId}

	err := safe.WaitCondition(20*time.Minute, 15*time.Second, func() (bool, error) {
		res, err := p.Cluster().GetVM(opts, m.Region)
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
	if err != nil {
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
				if ip4[0] == "192" || ip4[0] == "10" || ip4[0] == "172" {
					priipv4s = append(priipv4s, nic.IPaddress)
				} else {
					pubipv4s = append(pubipv4s, nic.IPaddress)
				}
			}
		}
	}

	ips[constants.PUBLICIPV4] = pubipv4s
	ips[constants.PRIVATEIPV4] = priipv4s
	return ips
}

func (m *Machine) mergeSameIPtype(mm map[string][]string) map[string][]string {
	for IPtype, ips := range mm {
		var sameIp string
		for _, ip := range ips {
			sameIp = sameIp + ip + ", "
		}
		if sameIp != "" {
			mm[IPtype] = []string{strings.TrimRight(sameIp, ", ")}
		}
	}
	return mm
}

func (m *Machine) Remove(p OneProvisioner) error {
	log.Debugf("  removing machine in one (%s)", m.Name)

	if m.VMId == "" {
		log.Debugf(" instance_id is empty removing machine in one (%s)", m.Name)
		return nil
	}

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

func (m *Machine) isDeleteOk() bool {
	return m.State != constants.StateInitialized && m.State != constants.StateInitializing && m.State != constants.StatePreError
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
	log.Debugf("  creating backup machine in one (%s)", m.Name)
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
	bk.ImageId = id
	return bk.UpdateBackup()
}

func (m *Machine) CreateDiskSnap(p OneProvisioner) error {
	log.Debugf("  creating snap machine in one (%s)", m.Name)
	snp, err := carton.GetSnap(m.CartonsId, m.AccountId)
	if err != nil {
		return err
	}
	vm := virtualmachine.Vnc{VmId: m.VMId}
	res, err := p.Cluster().GetVM(vm, m.Region)
	if err != nil {
		return err
	}

	expects := res.LenSnapshots() + 1
	vmid, _ := strconv.Atoi(m.VMId)
	opts := snapshot.Snapshot{
		VMId:            vmid,
		DiskId:          0,
		DiskDiscription: snp.Name,
	}

	id, err := p.Cluster().SnapVMDisk(opts, m.Region)
	if err != nil {
		return err
	}

	res, err = p.Cluster().GetVM(vm, m.Region)
	if err != nil {
		return err
	}
	if res.LenSnapshots() == expects {
		m.ImageId = id
		return nil
	}
	return fmt.Errorf(res.UserTemplate.Error)
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
		VMId:   vmid,
		DiskId: diskId,
		SnapId: sid,
	}
	m.ImageId = snp.SnapId
	return p.Cluster().RestoreSnap(opts, m.Region)
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

	return err
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
		if v.SnapId != m.ImageId && v.Status == constants.ACTIVESNAP {
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

func (m *Machine) UpdateBackupStatus(status utils.Status) error {
	bk, err := carton.GetBackup(m.CartonsId, m.AccountId)
	if err != nil {
		return err
	}
	bk.Status = status.String()
	return bk.UpdateBackup()
}

func (m *Machine) UpdateBackupPath(p OneProvisioner) error {
	bk, err := carton.GetBackup(m.CartonsId, m.AccountId)
	if err != nil {
		return err
	}
	var srcPath = make(map[string][]string)
	id, _ := strconv.Atoi(m.ImageId)
	opts := images.Image{
		Id: id,
	}
	res, err := p.Cluster().GetImage(opts, m.Region)
	if err != nil {
		return err
	}
	srcPath[constants.SOURCE_PATH] = []string{res.Source}
	srcPath[constants.DATASTORE_ID] = []string{strconv.Itoa(res.DatastoreID)}
	bk.Outputs.NukeAndSet(srcPath)
	bk.ImageId = m.ImageId
	return bk.UpdateBackup()
}

func (m *Machine) UpdateBackupVMIps() error {
	var ips = make(map[string][]string)
	asm, err := carton.NewAssembly(m.CartonId, m.AccountId, "")
	if err != nil {
		return err
	}

	bk, err := carton.GetBackup(m.CartonsId, m.AccountId)
	if err != nil {
		return err
	}

	nics := []string{constants.PUBLICIPV4, constants.PRIVATEIPV4, constants.PUBLICIPV6, constants.PRIVATEIPV6}
	for _, nic := range nics {
		ip := strings.Split(asm.Outputs.Match(nic), ",")
		if len(ip) > 0 && ip[0] != "" {
			ips[nic] = ip
		}
	}
	bk.Outputs.NukeAndSet(ips)
	return bk.UpdateBackup()
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
	log.Debugf("  remove backup image (%s) in one ", m.Name)
	opts := compute.Image{
		Name:    bk.Name,
		Region:  m.Region,
		ImageId: id,
	}
	err = p.Cluster().RemoveImage(opts)
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
		VMId:   vmid,
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
	if m.Status == constants.StatusDestroying {
		quota.AllocatedTo = ""
	}

	return quota.Update()
}

func (m *Machine) getImage(id string) (*mk.RawImages, error) {
	r := new(mk.RawImages)
	r.AccountId = m.AccountId
	r.Id = id
	return r.Get()
}

func (m *Machine) getMarketPlace(id string) (*mk.Marketplaces, error) {
	r := new(mk.Marketplaces)
	r.AccountId = m.AccountId
	r.Id = id
	return r.Get()
}

func (m *Machine) CreateImage(p OneProvisioner, img images.ImageType) error {
	opts := images.Image{
		Name: m.Name,
		Path: m.PublicUrl,
		Type: img,
	}

	res, err := p.Cluster().ImageCreate(opts, m.Region)
	if err != nil {
		return err
	}
	m.ImageId = res.(string)
	return nil
}

func (m *Machine) UpdateImage() error {
	raw, err := m.getImage(m.CartonId)
	if err != nil {
		return err
	}
	var id = make(map[string][]string)
	id[constants.RAW_IMAGE_ID] = []string{m.ImageId}
	raw.Status = string(m.Status)
	raw.Outputs.NukeAndSet(id)
	return raw.Update()
}

func (m *Machine) UpdateImageStatus() error {
	raw, err := m.getImage(m.CartonId)
	if err != nil {
		return err
	}
	raw.Status = string(m.Status)
	return raw.Update()
}

func (m *Machine) UpdateMarketImageId() error {
	mark, err := m.getMarketPlace(m.CartonId)
	if err != nil {
		return err
	}
	var id = make(map[string][]string)
	id[constants.IMAGE_ID] = []string{m.ImageId}
	mark.Outputs.NukeAndSet(id)
	return mark.Update()
}

func (m *Machine) UpdateMarketplaceStatus() error {
	mark, err := m.getMarketPlace(m.CartonId)
	if err != nil {
		return err
	}
	return mark.UpdateStatus(m.Status)
}

func (m *Machine) UpdateMarketplaceError(causeof error) error {
	mark, err := m.getMarketPlace(m.CartonId)
	if err != nil {
		return err
	}
	return mark.UpdateError(m.Status, causeof)
}

func (m *Machine) CreateDatablock(p OneProvisioner, box *provision.Box) error {
	size, _ := strconv.Atoi(strconv.FormatInt(int64(box.GetHDD()), 10))
	opts := images.Image{
		Name:       m.Name,
		Size:       size,
		Type:       images.DATABLOCK,
		Persistent: "yes",
	}
	res, err := p.Cluster().ImageCreate(opts, m.Region)
	if err != nil {
		return err
	}
	m.ImageId = res.(string)
	return nil
}

func (m *Machine) RemoveDatablock(p OneProvisioner) error {
	mark, err := m.getMarketPlace(m.CartonId)
	if err != nil {
		return err
	}
	id, _ := strconv.Atoi(mark.ImageId())
	opts := compute.Image{
		Region:  m.Region,
		ImageId: id,
	}
	err = p.Cluster().RemoveImage(opts)
	if err != nil {
		return err
	}
	return mark.NukeKeysOutputs(constants.IMAGE_ID)
}

func (m *Machine) CreateInstance(p OneProvisioner, box *provision.Box) error {
	var uname, rawname, imagename string

	mark, err := m.getMarketPlace(m.CartonId)
	if err != nil {
		return err
	}

	raw, err := m.getImage(mark.RawImageId())
	if err != nil {
		return err
	}
	rawname = raw.Name
	imagename = box.CartonName

	XMLtemplate, err := p.Cluster().GetTemplate(m.Region)
	if err != nil {
		return err
	}

	XMLtemplate.Template.Cpu = strconv.FormatInt(int64(box.GetCpushare()), 10)
	XMLtemplate.Template.VCpu = XMLtemplate.Template.Cpu
	XMLtemplate.Template.Memory = strconv.FormatInt(int64(box.GetMemory()), 10)
	XMLtemplate.Template.Cpu_cost = mark.GetVMCpuCost()
	XMLtemplate.Template.Memory_cost = mark.GetVMMemoryCost()
	XMLtemplate.Template.Disk_cost = mark.GetVMHDDCost()
	XMLtemplate.Template.Context.Accounts_id = box.AccountId
	XMLtemplate.Template.Context.Marketplace_id = box.CartonId
	XMLtemplate.Template.Context.ApiKey = box.ApiArgs.Api_Key
	XMLtemplate.Template.Context.Org_id = box.OrgId

	if len(XMLtemplate.Template.Disks) >= 0 {
		uname = XMLtemplate.Template.Disks[0].Image_Uname
	} else {
		uname = "oneadmin"
	}
	disks := make([]*template.Disk, 0)
	disks = append(disks, &template.Disk{Image_Uname: uname, Image: rawname})
	XMLtemplate.Template.Disks = disks
	vmid, err := p.Cluster().InstantiateVM(XMLtemplate, imagename, m.VCPUThrottle, m.Region)
	if err != nil {
		return err
	}

	m.VMId = vmid
	var id = make(map[string][]string)
	id[carton.INSTANCE_ID] = []string{m.VMId}
	if err = mark.NukeAndSetOutputs(id); err != nil {
		return err
	}
	return nil
}

func (m *Machine) AttachDatablock(p OneProvisioner, b *provision.Box) error {
	id, _ := strconv.Atoi(m.VMId)
	opts := &disk.VmDisk{
		VmId: id,
		Vm:   disk.Vm{Disk: disk.Disk{Image: b.CartonName}},
	}
	return p.Cluster().AttachDisk(opts, m.Region)
}

func (m *Machine) MarketplaceInstanceState(p OneProvisioner) error {
	opts := virtualmachine.Vnc{
		VmId: m.VMId,
	}

	mark, err := m.getMarketPlace(m.CartonId)
	if err != nil {
		return err
	}
	res := &virtualmachine.VM{}
	_ = mark.UpdateStatus(utils.Status(constants.StatusLcmStateChecking))

	err = safe.WaitCondition(30*time.Minute, 20*time.Second, func() (bool, error) {
		_ = mark.Trigger_event(utils.Status(constants.StatusWaitUntill))
		res, err = p.Cluster().GetVM(opts, m.Region)
		if err != nil {
			return false, err
		}
		if res.State == int(virtualmachine.DONE) {
			return false, fmt.Errorf("VM deleted while machine deploying")
		}
		status := res.StateString()
		if res.LcmStateString() != "" {
			status = status + "_" + res.LcmStateString()
		}
		_ = mark.Trigger_event(utils.Status(status))
		return (res.HistoryRecords.History != nil && res.LcmState == 3), nil
	})
	return err
}

func (m *Machine) GetMarketplaceVNC(p OneProvisioner) error {
	opts := virtualmachine.Vnc{
		VmId: m.VMId,
	}
	res, err := p.Cluster().GetVM(opts, m.Region)
	if err != nil {
		return err
	}
	//	ips := m.mergeSameIPtype(m.IPs(res.Nics()))
	m.VNCHost = res.GetHostIp()
	m.VNCPort = res.GetPort()
	return nil
}

func (m *Machine) UpdateMarketplaceVNC() error {
	var vnc = make(map[string][]string)
	vnc[carton.VNCHOST] = []string{m.VNCHost}
	vnc[carton.VNCPORT] = []string{m.VNCPort}

	if mark, err := m.getMarketPlace(m.CartonId); err != nil {
		return err
	} else if err = mark.NukeAndSetOutputs(vnc); err != nil {
		return err
	}
	return nil
}

func (m *Machine) ImagePersistent(p OneProvisioner) error {
	id, _ := strconv.Atoi(m.ImageId)
	opts := images.Image{
		Id: id,
	}
	return p.Cluster().ImagePersistent(opts, m.Region)
}

func (m *Machine) ImageTypeChange(p OneProvisioner) error {
	id, _ := strconv.Atoi(m.ImageId)
	opts := images.Image{
		Id:   id,
		Type: images.OPERATING_SYSTEM,
	}
	return p.Cluster().ImageTypeChange(opts, m.Region)
}

func (m *Machine) CheckSaveImage(p OneProvisioner) error {
	mark, err := m.getMarketPlace(m.CartonId)
	if err != nil {
		return err
	}
	m.ImageId = mark.ImageId()
	id, _ := strconv.Atoi(m.ImageId)
	opts := images.Image{
		Id: id,
	}
	res, err := p.Cluster().GetImage(opts, m.Region)
	if err != nil {
		return err
	}
	if res.Persistent == "no" {
		return fmt.Errorf("Image in Non-persistent state")
	}
	return m.WaitUntillVMState(p, virtualmachine.POWEROFF, virtualmachine.LCM_INIT)
}

func (m *Machine) StopMarkplaceInstance(p OneProvisioner) error {
	if res, err := p.Cluster().GetVM(virtualmachine.Vnc{VmId: m.VMId}, m.Region); err == nil {
		if res.State != int(virtualmachine.POWEROFF) {
			vmid, _ := strconv.Atoi(m.VMId)
			return p.Cluster().VM(compute.VirtualMachine{VMId: vmid, Region: m.Region}, "stop")
		}
	} else {
		return err
	}
	return nil
}

func (m *Machine) RemoveInstance(p OneProvisioner) error {
	mark, err := m.getMarketPlace(m.CartonId)
	if err != nil {
		return err
	}
	if mark.RemoveVM() == constants.YES {
		return m.Remove(p)
	} else {
		vmid, _ := strconv.Atoi(m.VMId)
		did, _ := strconv.Atoi("0")
		opts := &disk.VmDisk{
			VmId: vmid,
			Vm:   disk.Vm{Disk: disk.Disk{Disk_Id: did}},
		}
		return p.Cluster().DetachDisk(opts, m.Region)
	}
}

func (m *Machine) UpdatePolicyStatus(index int) error {
	asm, err := carton.NewAssembly(m.CartonId, m.AccountId, "")
	if err != nil {
		return err
	}
	return asm.UpdatePolicyStatus(index, m.Status)
}

func (m *Machine) AttachNetwork(box *provision.Box, p OneProvisioner) error {
	return p.Cluster().AttachNics(box.PolicyOps, m.VMId, m.Region, box.StorageType)
}

func (m *Machine) DetachNetwork(box *provision.Box, p OneProvisioner) error {
	ips := m.removableIPs(box.PolicyOps.Rules)
	res, err := p.Cluster().GetVM(virtualmachine.Vnc{VmId: m.VMId}, m.Region)
	if err != nil {
		return err
	}
	ids := m.networkIds(res, ips)
	return p.Cluster().DetachNics(ids, m.VMId, m.Region)
}

func (m *Machine) networkIds(vm *virtualmachine.VM, ips map[string]string) []string {
	var net_ids []string
	for _, ip := range ips {
		id := vm.NetworkIdByIP(ip)
		if id != "" {
			net_ids = append(net_ids, id)
		}
	}
	return net_ids
}

func (m *Machine) removableIPs(rules map[string]string) map[string]string {
	ips := make(map[string]string, 0)
	for _, key := range carton.NETWORK_KEYS {
		if ip, ok := rules[key]; ok {
			ips[key] = ip
		}
	}
	return ips
}

func (m *Machine) RemoveNetworkIps(box *provision.Box) error {
	asm, err := carton.NewAssembly(m.CartonId, m.AccountId, "")
	if err != nil {
		return err
	}
	mm := make(map[string][]string, 0)
	ips := m.removableIPs(box.PolicyOps.Rules)
	for key, ip := range ips {
		if value := asm.Outputs.Match(key); value != "" {
			mm[key] = []string{strings.Replace(value, ip, "", -1)}
		}
	}

	return asm.NukeAndSetOutputs(mm)
}

func (m *Machine) RemoveImage(p OneProvisioner) error {
	if m.ImageId == "" {
		return nil
	}
	id, _ := strconv.Atoi(m.ImageId)
	log.Debugf("  remove image in one (%s)", m.Name)
	opts := compute.Image{
		Region:  m.Region,
		ImageId: id,
	}
	return p.Cluster().RemoveImage(opts)
}

func (m *Machine) SetStatusErr(status utils.Status, causeof error) error {
	log.Debugf("  set status[%s] of machine (%s, %s)", m.Id, m.Name, status.String())

	if asm, err := carton.NewAssembly(m.CartonId, m.AccountId, ""); err != nil {
		return err
	} else if err = asm.SetStatusErr(status, causeof); err != nil {
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
