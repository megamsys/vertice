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
	"github.com/megamsys/megamd/db"
	"gopkg.in/yaml.v2"
	"strings"
)

//bunch Assemblys
type Cartons []*Carton

//The grand elephant for megam cloud platform.
type Assemblies struct {
	Id          string      `json:"id"`
	AccountsId  string      `json:"accounts_id"`
	JsonClaz    string      `json:"json_claz"`
	Name        string      `json:"name"`
	AssemblysId []string    `json:"assemblies"`
	Inputs      []*JsonPair `json:"inputs"`
	CreatedAt   string      `json:"created_at"`
}

func (a *Assemblies) String() string {
	if d, err := yaml.Marshal(a); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

/** A public function which pulls the assemblies for deployment.
and any others we do. **/
func Get(id string) (*Assemblies, error) {
	a := &Assemblies{}
	if err := db.Fetch("assemblies", id, a); err != nil {
		return nil, err
	}

	log.Debugf("Assemblies %v", a)
	return a, nil
}

//make cartons from assemblies.
func (a *Assemblies) MkCartons() (Cartons, error) {
	newCs := make(Cartons, 0, len(a.AssemblysId))
	for _, ay := range a.AssemblysId {
		if len(strings.TrimSpace(ay)) > 1 {
			if ca, err := mkCarton(a.Id, ay); err != nil {
				return nil, err
			} else {
				ca.toBox()                //on success, make a carton2box if BoxLevel is BoxZero
				newCs = append(newCs, ca) //on success append carton
			}
		}
	}
	log.Debugf("Cartons %v", newCs)
	return newCs, nil
}

//a hash in json representing {name: "", value: ""}
type JsonPairs []*JsonPair

type JsonPair struct {
	K string `json:"key"`
	V string `json:"value"`
}

//create a new hash pair in json  by providing a key, value
func NewJsonPair(k string, v string) *JsonPair {
	return &JsonPair{
		K: k,
		V: v,
	}
}

//match for a value in the JSONPair arrays and send the value
func (p *JsonPairs) match(k string) string {
	for _, j := range *p {
		if j.K == k {
			return j.V
		}
	}
	return ""
}
