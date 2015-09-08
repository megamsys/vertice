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
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/megamd/db"
	"github.com/megamsys/megamd/provision"
	"gopkg.in/yaml.v2"
)

var Provisioner provision.Provisioner

type Cartons []*Carton
type JsonPairs []*JsonPair

type JsonPair struct {
	K string `json:"key"`
	V string `json:"value"`
}

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

/** This is a plublic function as this is the  entry for the app deployment
and any others we do. **/
func Get(id string) (*Assemblies, error) {
	a := &Assemblies{}
	if err := db.Fetch("assemblies", id, a); err != nil {
		return nil, err
	}

	log.Debugf("Assemblies %v", a)
	return a, nil
}

func (a *Assemblies) MkCartons() (Cartons, error) {
	newCs := make(Cartons, 0, len(a.AssemblysId))
	for _, ai := range a.AssemblysId {
		if len(strings.TrimSpace(ai)) > 1 {
			if b, err := mkCarton(ai); err != nil {
				return nil, err
			} else {
				newCs = append(newCs, b)
			}
		}
	}
	log.Debugf("Cartons %v", newCs)
	return newCs, nil
}

//match for a value in the JSONPair and send the value
func (p *JsonPairs) match(k string) string {
	for _, j := range *p {
		if j.K == k {
			return j.V
		}
	}
	return ""
}
