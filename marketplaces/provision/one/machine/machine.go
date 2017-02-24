package machine

import (
	"github.com/megamsys/libgo/utils"
	// log "github.com/Sirupsen/logrus"
	 constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/opennebula-go/images"
	mk "github.com/megamsys/vertice/marketplaces"
	"github.com/megamsys/vertice/marketplaces/provision"
	"github.com/megamsys/vertice/marketplaces/provision/one/cluster"
  "strconv"
)

type OneProvisioner interface {
	Cluster() *cluster.Cluster
  Resource() map[string]string
}

type Machine struct {
	Name         string
	Region       string
	CartonId     string
	AccountId    string
	Image        string
	PublicUrl    string
	VCPUThrottle string
	VMId         string
	VNCHost      string
	VNCPort      string
	ImageId      string
	StorageType  string
	Status       utils.Status
}

type CreateArgs struct {
	Raw         *mk.RawImages
	Market      *mk.Marketplaces
	Provisioner OneProvisioner
}

func (m *Machine) getRaw() (*mk.RawImages, error) {
	r := new(mk.RawImages)
	r.AccountId = m.AccountId
	r.Id = m.CartonId
	return r.Get()
}

func (m *Machine) getMarket() (*mk.Marketplaces, error) {
	r := new(mk.Marketplaces)
	r.AccountId = m.AccountId
	r.Id = m.CartonId
	return r.Get()
}

func (m *Machine) CreateISO(p OneProvisioner) error {
	opts := images.Image{
		Name: m.Name,
		Path: m.PublicUrl,
		Type: images.CD_ROM,
	}

	res, err := p.Cluster().ImageCreate(opts, m.Region)
	if err != nil {
		return err
	}
	m.ImageId = res.(string)
	return nil
}

func (m *Machine) UpdateRawImageId() error {
  raw, err := m.getRaw()
  if err != nil {
    return err
  }
	var id = make(map[string][]string)
	id[constants.RAW_IMAGE_ID] = []string{m.ImageId}
  raw.Status = string(m.Status)
	raw.Outputs.NukeAndSet(id)
	return raw.Update()
}

func (m *Machine) UpdateRawStatus() error {
  raw, err := m.getRaw()
  if err != nil {
    return err
  }
  raw.Status = string(m.Status)
	return raw.Update()
}

func (m *Machine) IsImageReady(p OneProvisioner) error {
	id, _ := strconv.Atoi(m.ImageId)
	opts := &images.Image{
		Id: id,
	}
	return p.Cluster().IsImageReady(opts, m.Region)
}

func (m *Machine) UpdateMarketImageId() error {
  mark, err := m.getMarket()
  if err != nil {
    return err
  }
	var id = make(map[string][]string)
	id[constants.IMAGE_ID] = []string{m.ImageId}
	mark.Outputs.NukeAndSet(id)
	return mark.Update()
}

func (m *Machine) UpdateMarketStatus() error {
  mark, err := m.getMarket()
  if err != nil {
    return err
  }
	return mark.UpdateStatus(m.Status)
}

func (m *Machine) CreateDatablock(p OneProvisioner) error {
	mm := p.Resource()
	size, _ := strconv.Atoi(mm[constants.STORAGE])
	opts := images.Image{
		Name: m.Name,
		Size: size,
		Type: images.DATABLOCK,
	}
	res, err := p.Cluster().ImageCreate(opts, m.Region)
	if err != nil {
		return err
	}
	m.ImageId = res.(string)
	return nil
}

func (m *Machine) CreateInstance(p OneProvisioner,box *provision.Box)  error {
	// resources := p.Resource()
	// opts := compute.VirtualMachine{
	// 	Name:       box.Name,
	// 	Image:      m.ImageName,
	// 	Region:     args.Box.Region,
	// 	Cpu:        strconv.FormatInt(int64(args.Box.GetCpushare()), 10),
	// 	Memory:     strconv.FormatInt(int64(args.Box.GetMemory()), 10),
	// 	HDD:        strconv.FormatInt(int64(args.Box.GetHDD()), 10),
	// 	CpuCost:    asm.GetVMCpuCost(),
	// 	MemoryCost: asm.GetVMMemoryCost(),
	// 	HDDCost:    asm.GetVMHDDCost(),
	// 	ContextMap: map[string]string{compute.ASSEMBLY_ID: args.Box.CartonId, compute.ORG_ID: args.Box.OrgId,
	// 		compute.ASSEMBLIES_ID: args.Box.CartonsId, compute.ACCOUNTS_ID: args.Box.AccountId, compute.API_KEY: args.Box.ApiArgs.Api_Key, constants.QUOTA_ID: args.Box.QuotaId},
	// 	Vnets: args.Box.Vnets,
	// }
	// opts.VCpu = opts.Cpu
	// if strings.Contains(args.Box.Tosca, "freebsd") {
	// 	opts.Files = "/detio/freebsd/init.sh"
	// }
	//
	// _, _, vmid, err := args.Provisioner.Cluster().CreateVM(opts, m.VCPUThrottle, m.StorageType)
	// if err != nil {
	// 	return err
	// }
	// m.VMId = vmid
	// var id = make(map[string][]string)
	// id[carton.INSTANCE_ID] = []string{m.VMId}
	// if err = asm.NukeAndSetOutputs(id); err != nil {
	// 	return err
	// }
	return nil
}
