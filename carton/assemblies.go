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
	log "github.com/golang/glog"
	"github.com/megamsys/libgo/db"
)

var Provisioner provision.Provisioner

type Cartons []*Carton

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

/** This is a plublic function as this is the  entry for the app deployment
and any others we do. **/
func Get(id string) (*Assemblies, error) {
	a := &Assemblies{}

	log.Infof("Assemblies %s", id)
	if conn, err := db.Conn("assemblies"); err != nil {
		return a, err
	}

	if err = conn.FetchStruct(id, a); err != nil {
		return a, err
	}
	defer conn.Close()
	log.Infof("Assemblies %v", a)
	return a, nil
}

func (a *Assemblies) mkCartons() (Cartons, error) {
	newCs := make(Cartons, len(a.AssemblysId))
	err := nil
	for i, ai := range a.AssemblysId {
		if err, b := carton.mkCarton(ai); err != nil {
			append(newCs, b)
		}
	}
	return newCs, err
}
