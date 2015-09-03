package carton

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/megamsys/megamd/action"
	"github.com/megamsys/megamd/caton/bind"
	"github.com/megamsys/megamd/db"
	"github.com/megamsys/megamd/provision"
	"github.com/megamsys/megamd/repository"
)

var (
	cnameRegexp = regexp.MustCompile(`^(\*\.)?[a-zA-Z0-9][\w-.]+$`)
)

type Carton struct {
	Name     string
	Tosca    string
	Cpushare string
	Memory   string
	Swap     string
	HDD      string
	Envs     map[string]bind.EnvVar
	Boxes    *[]provision.Box
}

// CartonDeploys.. we need the config or no ?
// May be stick it into the provisioner directly during a load.
func CDeploys(c *Carton) {
	for box := range c.Boxes {
		err := carton.Deploys(&DeployOpts{B: box})
		if err != nil {
			log.Errorf("Unable to destroy box in provisioner", err)
		}
	}
}

// Deletes a carton.
func Delete(c *Carton) {
	for box := range c.Boxes {
		err := Provisioner.Destroy(box)
		if err != nil {
			log.Errorf("Unable to destroy box in provisioner", err)
		}
	}
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
func (c *Carton) Group() ([]*YBoundBox, err) {
	return nil, nil
}

// Available returns true if at least one of N units is started or unreachable.
func (c *Carton) Available() bool {
	boxes, err := c.Boxes()
	if err != nil {
		return false
	}
	for _, box := range boxes {
		if box.Available() {
			return true
		}
	}
	return false
}

// Start starts the app calling the provisioner.Start method and
// changing the units state to StatusStarted.
func (c *Carton) Start() error {
	for box := range c.Boxes {
		err := Provisioner.Start(box)
		if err != nil {
			log.Errorf("[start] error on start the box %s - %s", box.Name, err)
			return err
		}
	}
	return nil
}

func (c *Carton) Stop() error {
	for box := range c.Boxes {
		err := Provisioner.Stop(box)
		if err != nil {
			log.Errorf("[start] error on start the box %s - %s", box.Name, err)
			return err
		}
	}
	return nil
}

// Restart runs the restart hook for the app, writing its output to w.
func (c *Carton) Restart() error {
	for box := range c.Boxes {
		err := Provisioner.Restart(box)
		if err != nil {
			log.Errorf("[start] error on start the box %s - %s", box.Name, err)
			return err
		}
	}
	return nil
}

// GetTosca returns the tosca type  of the carton.
func (c *Carton) GetTosca() string {
	return c.Tosca
}

// GetMemory returns the memory limit (in bytes) for the carton.
func (c *Carton) GetMemory() int64 {
	return c.Memory
}

// GetCpuShare returns the cpu share for the carton.
func (c *Carton) GetCpuShare() int {
	return c.CpuShare
}

// GetMemory returns the memory limit (in bytes) for the carton.
func (c *Carton) GetHDD() int64 {
	return c.HDD
}

// Envs returns a map representing the apps environment variables.
func (c *Carton) Envs() map[string]bind.EnvVar {
	return c.Env
}

// AddCName adds a CName to box. It updates the attribute,
// calls the SetCName function on the provisioner and saves
// the box in the database, returning an error when it cannot save the change
// in the database or add the CName on the provisioner.
func (b *Box) AddCName(cnames ...string) error {
	for _, cname := range cnames {
		if cname != "" && !cnameRegexp.MatchString(cname) {
			return stderr.New("Invalid cname")
		}

		if s, ok := Provisioner.(provision.CNameManager); ok {
			if err := s.SetCName(app, cname); err != nil {
				return err
			}
		}
		//Riak: append the ip/cname in the component.
		//here (or) can be handled as an action.
	}
	return nil
}

func (c *Carton) RemoveCName(cnames ...string) error {
	for _, cname := range cnames {
		count := 0
		for _, appCname := range app.CName {
			if cname == appCname {
				count += 1
			}
		}
		if count == 0 {
			return stderr.New("cname not exists!")
		}
		if s, ok := Provisioner.(provision.CNameManager); ok {
			if err := s.UnsetCName(app, cname); err != nil {
				return err
			}
		}
		//Riak: append the ip/cname in the component available in the box.
		//or handle it as an action
		if err != nil {
			return err
		}
	}
	return nil
}

// Log adds a log message to the app. Specifying a good source is good so the
// user can filter where the message come from.
func (box *Box) Log(message, source, unit string) error {
	messages := strings.Split(message, "\n")
	logs := make([]interface{}, 0, len(messages))
	for _, msg := range messages {
		if msg != "" {
			l := Applog{
				Date:    time.Now().In(time.UTC),
				Message: msg,
				Source:  source,
				AppName: app.Name,
				Unit:    unit,
			}
			logs = append(logs, l)
		}
	}
	if len(logs) > 0 {
		notify(app.Name, logs)
		conn, err := db.LogConn()
		if err != nil {
			return err
		}
		defer conn.Close()
		return conn.Logs(app.Name).Insert(logs...)
	}
	return nil
}
