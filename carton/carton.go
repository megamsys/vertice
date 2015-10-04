package carton

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/megamd/carton/bind"
	"github.com/megamsys/megamd/provision"
	"github.com/megamsys/megamd/repository"
	"gopkg.in/yaml.v2"
)

// A carton represents a real world assembly.
// This struct provides and easy way to manage information about an assembly, instead passing it around
type Carton struct {
	Id           string //assemblyid
	Name         string
	CartonsId    string
	Tosca        string
	ImageVersion string
	Compute      provision.BoxCompute
	Repo         repository.Repo
	DomainName   string
	Provider     string
	PublicIp     string
	Envs         []bind.EnvVar
	Boxes        *[]provision.Box
}

func (a *Carton) String() string {
	if d, err := yaml.Marshal(a); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

//A global provisioner set by the subd daemon.
//A BUG, the Provisioner can't be a global variable as
//this will be overwritten if multiple subd daemons set something.
var Provisioner provision.Provisioner

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
			CartonId:     c.Id,        //this isn't needed.
			Id:           c.Id,        //assembly id sent in ContextMap
			CartonsId:    c.CartonsId, //assembliesId,
			Level:        c.lvl(),     //based on the level, we decide to use the Box-Id as ComponentId or AssemblyId
			Name:         c.Name,
			ImageVersion: c.ImageVersion,
			DomainName:   c.DomainName,
			Compute:      c.Compute,
			Repo:         c.Repo,
			Provider:     c.Provider,
			PublicIp:     c.PublicIp,
			Tosca:        c.Tosca,
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
		err := ChangeState(&StateChangeOpts{B: &box, Changed: provision.StatusStateup})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Carton) Bind(box *provision.Box) error {
	/*	bis := c.Group()
		for _, bi := range bis {
			err = bi.Bind(c, bi)
			if err != nil {
				log.Errorf("Error binding the box %s with the service %s: %s", box.Name, bi.Name, err)
			}
		}
	*/
	return nil
}

func (c *Carton) Unbind(box *provision.Box) error {
	/*	bis := c.Group()
		for _, bi := range bis {
			err = bi.UnBind(c, bi)
			if err != nil {
				log.Errorf("Error binding the box %s with the service %s: %s", box.Name, bi.Name, err)
			}
		}*/
	return nil
}

// Group the related_components into BindInstances.
func (c *Carton) Group() ([]*bind.YBoundBox, error) {
	return nil, nil
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

// starts the box calling the provisioner.
// changing the boxes state to StatusStarted.
func (c *Carton) Start() error {
	for _, box := range *c.Boxes {
		err := Provisioner.Start(&box, "", nil)
		if err != nil {
			log.Errorf("Unable to start the box  %s", err)
			return err
		}
	}
	return nil
}

// stops the box calling the provisioner.
// changing the boxes state to StatusStopped.
func (c *Carton) Stop() error {
	for _, box := range *c.Boxes {
		err := Provisioner.Stop(&box, "", nil)
		if err != nil {
			log.Errorf("Unable to stop the box %s", err)
			return err
		}
	}
	return nil
}

// restarts the box calling the provisioner.
// changing the boxes state to StatusStarted.
func (c *Carton) Restart() error {
	for _, box := range *c.Boxes {
		err := Provisioner.Restart(&box, "", nil)
		if err != nil {
			log.Errorf("[start] error on start the box %s - %s", box.Name, err)
			return err
		}
	}
	return nil
}

// Envs returns a map representing the apps environment variables.
func (c *Carton) GetEnvs() []bind.EnvVar {
	return c.Envs
}
