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
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/db"
	"github.com/megamsys/megamd/global"
	"github.com/megamsys/megamd/log"
	"github.com/megamsys/megamd/provisioner"
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
	Requirements []*JsonPair   `json:"requirements"`
	Policies     []*Policy     `json:"policies"`
	Inputs       []*JsonPair   `json:"inputs"`
	Operations   []*Operations `json:"operations"`
	Outputs      []*JsonPair   `json:"outputs"`
	Status       string        `json:"status"`
	CreatedAt    string        `json:"created_at"`
}

type Assembly struct {
	ambly
	ComponentIds []string `json:"components"`
	Components   map[string]Component
}

//mkAssemblies into a carton. Just use what you need inside this carton
//a carton comprises of self contained boxes (actually a "colored component") externalized
//with what we need.
func GetCarton(id string) (*Carton, error) {
	a, err := get(id)
	if err != nil {
		return nil, err
	}

	c = &Carton{
		Name:     a.Name,
		Platform: a.Tosca,
		Cpushare: a.Cpushare,
		Memory:   a.Memory,
		Swap:     a.Swap,
		HDD:      a.HDD,
		Envs:     a.envs(),
		Boxes:    a.mkBoxes(),
	}
	return c, nil
}

//get the assebmly and its full detail of a component. we only store the
//componentid, hence you see that we have a components map to cater to that need.
func get(id string) (*Assembly, error) {
	log.Infof("get %s", id)
	a := &Assembly{}
	if conn, err := db.Conn("assembly"); err != nil {
		return d, err
	}

	if err := conn.FetchStruct(id, a); err != nil {
		return d, ferr
	}
	a.dig()
	defer conn.Close()
	return result, nil
}

func (a *Assembly) dig() error {
	for i, cid := range a.ComponentIds {
		if len(strings.TrimSpace(cid)) > 1 {
			comp := NewComponent(cid)
			if err := comp.Get(comp.Id); err != nil {
				log.Errorf("Failed to get component %s from riak: %s.", comp.Id, err.Error())
				return err
			}
			a.Components[comp.Id] = comp
		}
	}
}

//lets make boxes with components to be mutated later or, and the required
//information for a launch.
func (a *Assembly) mkBoxes() ([]*provision.Box, error) {
	err := nil
	newBoxs := make([]*provision.Box, len(a.Components))

	for i, comp := range a.Components {
		if err, b := comp.mkBox(); err != nil {
			append(newBoxs, b)
		}
	}
	return newBoxs, err
}

//all the variables in the inputs shall be treated as ENV.
//we can use a filtered approach as well.
func (a *Assembly) envs() {
	envs := make([]bind.EnvVar, 0, len(a.Inputs))
	for k, v := range a.Inputs {
		envs = append(envs, bind.EnvVar{ Name: k, Value: v})
	}
}
