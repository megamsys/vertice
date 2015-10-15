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
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/megamd/carton/bind"
	"github.com/megamsys/megamd/db"
	"github.com/megamsys/megamd/provision"
	"gopkg.in/yaml.v2"
)

type Policy struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Members []string `json:"members"`
}

//An assembly comprises of various components.
type ambly struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	JsonClaz  string    `json:"json_claz"`
	Tosca     string    `json:"tosca_type"`
	Inputs    JsonPairs `json:"inputs"`
	Outputs   JsonPairs `json:"outputs"`
	Policies  []*Policy `json:"policies"`
	Status    string    `json:"status"`
	CreatedAt string    `json:"created_at"`
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

//Assembly into a carton.
//a carton comprises of self contained boxes
func mkCarton(aies string, ay string) (*Carton, error) {
	a, err := get(ay)
	if err != nil {
		return nil, err
	}

	b, err := a.mkBoxes(aies)
	if err != nil {
		return nil, err
	}

	c := &Carton{
		Id:           ay,   //assembly id
		CartonsId:    aies, //assemblies id
		Name:         a.Name,
		Tosca:        a.Tosca,
		ImageVersion: a.imageVersion(),
		Envs:         a.envs(),
		DomainName:   a.domain(),
		Compute:      a.newCompute(),
		Provider:     a.provider(),
		PublicIp:     a.publicIp(),
		Boxes:        &b,
	}
	return c, nil
}

//lets make boxes with components to be mutated later or, and the required
//information for a launch.
//A "colored component" externalized with what we need.
func (a *Assembly) mkBoxes(aies string) ([]provision.Box, error) {
	newBoxs := make([]provision.Box, 0, len(a.Components))

	for _, comp := range a.Components {
		if len(strings.TrimSpace(comp.Id)) > 1 {
			if b, err := comp.mkBox(); err != nil {
				return nil, err
			} else {
				b.CartonId = a.Id
				b.CartonsId = aies
				b.CartonName = a.Name
				if b.Repo.IsEnabled() {
					b.Repo.Hook.CartonId = a.Id //this is screwy, why do we need it.
					b.Repo.Hook.BoxId = comp.Id
				}
				b.Compute = a.newCompute()
				newBoxs = append(newBoxs, b)
			}
		}
	}
	return newBoxs, nil
}

//Temporary hack to create an assembly from its id.
//This is used by SetStatus.
//We need add a Notifier interface duck typed by Box and Carton ?
func NewAssembly(id string) (*Assembly, error) {
	return get(id)
}

func (a *Assembly) SetStatus(status provision.Status) error {
	LastStatusUpdate := time.Now().In(time.UTC)

	a.Inputs = append(a.Inputs, NewJsonPair("lastsuccessstatusupdate", LastStatusUpdate.String()))
	a.Inputs = append(a.Inputs, NewJsonPair("status", status.String()))

	if err := db.Store(BUCKET, a.Id, a); err != nil {
		return err
	}
	return nil

}

//get the assebmly and its children (component). we only store the
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

func (a *Assembly) domain() string {
	return a.Inputs.match(DOMAIN)
}

func (a *Assembly) provider() string {
	return a.Inputs.match(provision.PROVIDER)
}

func (a *Assembly) publicIp() string {
	return a.Outputs.match(PUBLICIP)
}

func (a *Assembly) imageVersion() string {
	return a.Inputs.match(IMAGE_VERSION)
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

func (a *Assembly) newCompute() provision.BoxCompute {
	return provision.BoxCompute{
		Cpushare: a.getCpushare(),
		Memory:   a.getMemory(),
		Swap:     a.getSwap(),
		HDD:      a.getHDD(),
	}
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

//The default HDD is 10. we should configure it in the megamd.conf
func (a *Assembly) getHDD() string {
	if len(strings.TrimSpace(a.Inputs.match(provision.HDD))) <= 0 {
		return "10"
	}
	return a.Inputs.match(provision.HDD)
}
