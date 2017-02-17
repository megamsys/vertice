package carton

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/api"
	"github.com/megamsys/libgo/utils"
	"github.com/megamsys/vertice/provision"
	"gopkg.in/yaml.v2"
)

// A carton represents a real world assembly.
// This struct provides and easy way to manage information about an assembly, instead passing it around
type Carton struct {
	Id           string //assemblyid
	Name         string
	CartonsId    string
	AccountId    string
	QuotaId      string
	ApiArgs      api.ApiArgs
	OrgId        string
	Tosca        string
	ImageVersion string
	ImageName    string
	StorageType  string
	Backup       bool
	Compute      provision.BoxCompute
	SSH          provision.BoxSSH
	DomainName   string
	Provider     string
	PublicIp     string
	InstanceId   string
	Region       string
	Vnets        map[string]string
	Boxes        *[]provision.Box
	Status       utils.Status
	State        utils.State
}


//Global provisioners set by the subd daemons.
var ProvisionerMap map[string]provision.Provisioner = make(map[string]provision.Provisioner)

func (a *Carton) String() string {
	if d, err := yaml.Marshal(a); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

//If there are boxes, then it set the enum BoxSome or its BoxZero
func (c *Carton) lvl() provision.BoxLevel {
	if len(*c.Boxes) > 0 {
		return provision.BoxSome
	} else {
		return provision.BoxNone
	}
}

//Converts a carton to a box, (applicable in torpedo case)
func (c *Carton) toBox() error { //assemblies id.
	switch c.lvl() {
	case provision.BoxNone:
		c.Boxes = &[]provision.Box{provision.Box{
			Id:           c.Id, //should be the component id, but in case of BoxNone there is no component id.
			AccountId:    c.AccountId,
			CartonId:     c.Id,        //We stick the assemlyid here.
			CartonsId:    c.CartonsId, //assembliesId,
			OrgId:        c.OrgId,
			ApiArgs:      c.ApiArgs,
			CartonName:   c.Name,
			Name:         c.Name,
			DomainName:   c.DomainName,
			StorageType:  c.StorageType,
			Level:        c.lvl(), //based on the level, we decide to use the Box-Id as ComponentId or AssemblyId
			ImageVersion: c.ImageVersion,
			ImageName:    c.ImageName,
			Backup:       c.Backup,
			Compute:      c.Compute,
			Provider:     c.Provider,
			PublicIp:     c.PublicIp,
			InstanceId:   c.InstanceId,
			QuotaId:      c.QuotaId,
			Region:       c.Region,
			Vnets:        c.Vnets,
			Tosca:        c.Tosca,
			Status:       c.Status,
			State:        c.State,
		},
		}
	}
	return nil
}

// Deploy carton, which basically deploys the boxes.
func (c *Carton) Deploy() error {
	for _, box := range *c.Boxes {
		err := Deploy(&DeployOpts{B: &box})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Carton) Running() error {
	for _, box := range *c.Boxes {
		err := Running(&DeployOpts{B: &box})
		if err != nil {
			return err
		}
	}
	return nil
}

// Destroys a carton, which deletes its boxes.
func (c *Carton) Destroy() error {
	for _, box := range *c.Boxes {
		err := Destroy(&DestroyOpts{B: &box})
		if err != nil {
			return err
		}
	}
	return nil
}

// moves the state to the desired state
// changing the boxes state to StatusStateup.
func (c *Carton) Stateup() error {
	for _, box := range *c.Boxes {
		err := ChangeState(&StateChangeOpts{B: &box, Changed: utils.StatusStateupped})
		if err != nil {
			return err
		}
	}
	return nil
}

// Available returns true if at least one of N boxes which is started
func (c *Carton) Available() bool {
	for _, box := range *c.Boxes {
		if box.Available() {
			return true
		}
	}
	return false
}

//upgrade run thru all the ops.
func (c *Carton) Upgrade() error {
	for _, box := range *c.Boxes {
		err := NewUpgradeable(&box).Upgrade()
		if err != nil {
			log.Errorf("Unable to upgrade box : %s", err)
			return err
		}
	}
	return nil
}

// starts box
func (c *Carton) Start() error {
	for _, box := range *c.Boxes {
		err := Start(&LifecycleOpts{B: &box})
		if err != nil {
			log.Errorf("Unable to start the box  %s", err)
			return err
		}
	}
	return nil
}

// stops the box
func (c *Carton) Stop() error {
	for _, box := range *c.Boxes {
		err := Stop(&LifecycleOpts{B: &box})
		if err != nil {
			log.Errorf("Unable to stop the box %s", err)
			return err
		}
	}
	return nil
}

// restarts the box
func (c *Carton) Restart() error {
	for _, box := range *c.Boxes {
		err := Restart(&LifecycleOpts{B: &box})
		if err != nil {
			log.Errorf("Unable to restart the box %s", err)
			return err
		}
	}
	return nil
}

// SnapCreate a carton, which creates an image by current state of its box.
func (c *Carton) SaveImage() error {
	for _, box := range *c.Boxes {
		err := SaveImage(&DiskOpts{B: &box})
		if err != nil {
			return err
		}
	}
	return nil
}

// SnapDelete a carton, which removes an existing image created from state of its box.
func (c *Carton) DeleteImage() error {
	for _, box := range *c.Boxes {
		err := DeleteImage(&DiskOpts{B: &box})
		if err != nil {
			return err
		}
	}
	return nil
}

// SnapCreate a carton, which creates an image by current state of its box.
func (c *Carton) CreateSnapshot() error {
	for _, box := range *c.Boxes {
		err := CreateSnapshot(&DiskOpts{B: &box})
		if err != nil {
			return err
		}
	}
	return nil
}


// SnapCreate a carton, which creates an image by current state of its box.
func (c *Carton) SnapshotSaveAs() error {
	for _, box := range *c.Boxes {
		err := SnapshotSaveAs(&DiskOpts{B: &box})
		if err != nil {
			return err
		}
	}
	return nil
}


// SnapDelete a carton, which removes an existing image created from state of its box.
func (c *Carton) DeleteSnapshot() error {
	for _, box := range *c.Boxes {
		err := DeleteSnapshot(&DiskOpts{B: &box})
		if err != nil {
			return err
		}
	}
	return nil
}

// SnapDelete a carton, which removes an existing image created from state of its box.
func (c *Carton) RestoreSnapshot() error {
	for _, box := range *c.Boxes {
		err := RestoreSnapshot(&DiskOpts{B: &box})
		if err != nil {
			return err
		}
	}
	return nil
}

// AttachDisk a carton, which creates a disk storage by current state of its box.
func (c *Carton) AttachDisk() error {
	for _, box := range *c.Boxes {
		err := AttachDisk(&DiskOpts{B: &box})
		if err != nil {
			return err
		}
	}
	return nil
}

// DetachDisk a carton, which removes an existing disk storage by current state of its box.
func (c *Carton) DetachDisk() error {
	for _, box := range *c.Boxes {
		err := DetachDisk(&DiskOpts{B: &box})
		if err != nil {
			return err
		}
	}
	return nil
}
