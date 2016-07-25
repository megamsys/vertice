package carton

import (
	log "github.com/Sirupsen/logrus"
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
	AccountsId   string
	Tosca        string
	ImageVersion string
	Compute      provision.BoxCompute
	SSH          provision.BoxSSH
	DomainName   string
	Provider     string
	PublicIp     string
	VMId         string
	Region       string
	Vnets        map[string]string
	Boxes        *[]provision.Box
	Status       utils.Status
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
			AccountsId:   c.AccountsId,
			CartonId:     c.Id,        //We stick the assemlyid here.
			CartonsId:    c.CartonsId, //assembliesId,
			CartonName:   c.Name,
			Name:         c.Name,
			DomainName:   c.DomainName,
			Level:        c.lvl(), //based on the level, we decide to use the Box-Id as ComponentId or AssemblyId
			ImageVersion: c.ImageVersion,
			Compute:      c.Compute,
			Provider:     c.Provider,
			PublicIp:     c.PublicIp,
			VMId:         c.VMId,
			Region:       c.Region,
			Vnets:        c.Vnets,
			Tosca:        c.Tosca,
			Status:       c.Status,
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

// DiskSave a carton, which creates an image by current state of its box.
func (c *Carton) SaveImage() error {
	for _, box := range *c.Boxes {
		err := SaveImage(&DiskSaveOpts{B: &box})
		if err != nil {
			return err
		}
	}
	return nil
}
