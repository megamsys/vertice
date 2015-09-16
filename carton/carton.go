package carton

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/megamd/carton/bind"
	"github.com/megamsys/megamd/provision"
	"gopkg.in/yaml.v2"
)

type BoxLevel int

const (
	// BoxAny indicates that there is atleast one box to deploy or delete.
	BoxAny BoxLevel = iota

	// BoxZero indicates that there are no boxes to deploy or delete.
	BoxZero
)

type Carton struct {
	Name       string
	AssemblyId string
	Tosca      string
	Envs       []bind.EnvVar
	Boxes      *[]provision.Box
}

func (a *Carton) String() string {
	if d, err := yaml.Marshal(a); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

//If there are boxes, then it set the enum BoxAny or its BoxZero
func (c *Carton) lvl() BoxLevel {
	if len(*c.Boxes) > 0 {
		return BoxAny
	} else {
		return BoxZero
	}
}

//Converts a carton to a box, if there are no boxes below.
func (c *Carton) toBox() error {
	switch c.lvl() {
	case BoxZero:
		c.Boxes = &[]provision.Box{provision.Box{
			AssemblyId: c.AssemblyId,
			Name:       c.Name,
			DomainName: "",
			Tosca:      c.Tosca,
		},
		}
	}
	return nil
}

// Deploy carton, which basically deploys the boxes.
func (c *Carton) Create() error {
	for _, box := range *c.Boxes {
		err := Create(&DeployOpts{B: &box})
		if err != nil {
			log.Errorf("Unable to create box", err)
		}
	}
	return nil
}

// Deletes a carton, which deletes its boxes.
func (c *Carton) Delete() error {
	for _, box := range *c.Boxes {
		err := Provisioner.Destroy(&box, nil)
		if err != nil {
			log.Errorf("Unable to destroy box", err)
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
		err := Provisioner.Start(&box, "")
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
		err := Provisioner.Stop(&box, "")
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

// moves the state to the desired state
// changing the boxes state to StatusStateup.
func (c *Carton) Stateup() error {
	for _, box := range *c.Boxes {
		err := Deploy(&DeployOpts{B: &box})
		if err != nil {
			log.Errorf("Unable to deploy box", err)
		}
	}
}

// moves the state down to the desired state
// changing the boxes state to StatusStatedown.
func (c *Carton) Statedown() error {
	return nil
}

// GetTosca returns the tosca type  of the carton.
func (c *Carton) GetTosca() string {
	return c.Tosca
}


// Envs returns a map representing the apps environment variables.
func (c *Carton) GetEnvs() []bind.EnvVar {
	return c.Envs
}
