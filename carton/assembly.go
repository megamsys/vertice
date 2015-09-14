/*
** Copyright [2013-2015] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */
package carton

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/megamd/carton/bind"
	"github.com/megamsys/megamd/db"
	"github.com/megamsys/megamd/provision"
	"gopkg.in/yaml.v2"
	"strings"
)

// Carton is the main type in megam. A carton represents a real world assembly.
// An assembly comprises of various components.
// This struct provides and easy way to manage information about an assembly, instead passing it around

type Policy struct {
	Name    string   `json:"name"`
	Ptype   string   `json:"ptype"`
	Members []string `json:"members"`
}

type ambly struct {
	Id           string        `json:"id"`
	Name         string        `json:"name"`
	JsonClaz     string        `json:"json_claz"`
	ToscaType    string        `json:"tosca_type"`
	Requirements JsonPairs     `json:"requirements"`
	Policies     []*Policy     `json:"policies"`
	Inputs       JsonPairs     `json:"inputs"`
	Operations   []*Operations `json:"operations"`
	Outputs      JsonPairs     `json:"outputs"`
	Status       string        `json:"status"`
	CreatedAt    string        `json:"created_at"`
}

type Assembly struct {
	ambly
	ComponentIds []string `json:"components"`
	Components   map[string]*Component
}

func (a *Assembly) String() string {
	if d, err := yaml.Marshal(a); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

func (a *Assembly) NewCompute() provision.BoxCompute {
	return provision.BoxCompute{
		Cpushare: a.getCpushare(),
		Memory:   a.getMemory(),
		Swap:     a.getSwap(),
		HDD:      a.getHDD(),
	}
}

//mkAssemblies into a carton. Just use what you need inside this carton
//a carton comprises of self contained boxes (actually a "colored component") externalized
//with what we need.
func mkCarton(id string) (*Carton, error) {
	a, err := get(id)
	if err != nil {
		return nil, err
	}

	b, err := a.mkBoxes()
	if err != nil {
		return nil, err
	}

	c := &Carton{
		AssemblyId: id,
		Name:       a.Name,
		Tosca:      a.ToscaType,
		Envs:       a.envs(),
		Boxes:      &b,
	}

	return c, nil
}

//get the assebmly and its full detail of a component. we only store the
//componentid, hence you see that we have a components map to cater to that need.
func get(id string) (*Assembly, error) {
	a := &Assembly{Components: make(map[string]*Component)}
	if err := db.Fetch("assembly", id, a); err != nil {
		return nil, err
	}
	a.dig()
	return a, nil
}

func (a *Assembly) dig() error {
	for _, cid := range a.ComponentIds {
		if len(strings.TrimSpace(cid)) > 1 {
			if comp, err := NewComponent(cid); err != nil {
				log.Errorf("Failed to get component %s from riak: %s.", cid, err.Error())
				return err
			} else {
				a.Components[cid] = comp
			}
		}
	}
	return nil
}

//lets make boxes with components to be mutated later or, and the required
//information for a launch.
func (a *Assembly) mkBoxes() ([]provision.Box, error) {
	newBoxs := make([]provision.Box, 0, len(a.Components))

	for _, comp := range a.Components {
		if len(strings.TrimSpace(comp.Id)) > 1 {
			if b, err := comp.mkBox(); err != nil {
				return nil, err
			} else {
				b.AssemblyId = a.Id
				b.Compute = a.NewCompute()
				newBoxs = append(newBoxs, b)
			}
		}
	}
	return newBoxs, nil
}

//all the variables in the inputs shall be treated as ENV.
//we can use a filtered approach as well.
func (a *Assembly) envs() []bind.EnvVar {
	envs := make([]bind.EnvVar, 0, len(a.Inputs))
	for _, i := range a.Inputs {
		envs = append(envs, bind.EnvVar{Name: i.K, Value: i.V})
	}
	return envs
}

func (a *Assembly) getCpushare() string {
	return a.Inputs.match(provision.CPU)
}

func (a *Assembly) getMemory() string {
	return a.Inputs.match(provision.RAM)

}

func (a *Assembly) getSwap() string {
	return ""
}

//The default HDD is 10.
func (a *Assembly) getHDD() string {
	return a.Inputs.match(provision.HDD)
}
