/*
** Copyright [2013-2016] [Megam Systems]
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
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	ldb "github.com/megamsys/libgo/db"
	"github.com/megamsys/vertice/meta"
	"strings"
)

type Payload struct {
	Id        string `json:"id" cql:"id"`
	Action    string `json:"action" cql:"action"`
	CatId     string `json:"cat_id" cql:"cat_id"`
	CatType   string `json:"cattype" cql:"cattype"`
	Category  string `json:"category" cql:"category"`
	CreatedAt string `json:"created_at" cql:"created_at"`
}

type PayloadConvertor interface {
	Convert(p *Payload) (*Requests, error)
}

func NewPayload(b []byte) (*Payload, error) {
	p := &Payload{}
	err := json.Unmarshal(b, &p)
	if err != nil {
		return nil, err
	}
	return p, err
}

/**
**fetch the request json from riak and parse the json to struct
**/
func (p *Payload) Convert() (*Requests, error) {
	if len(strings.TrimSpace(p.Id)) > 1 {
		return listReqsById(p.Id)
	} else {
		return &Requests{
			Action:    p.Action,
			Category:  p.Category,
			CatId:     p.CatId,
			CreatedAt: p.CreatedAt,
		}, nil
	}

}

//The payload in the queue can be just a pointer or a value.
//pointer means just the id will be available and rest is blank.
//value means the id is blank and others are available.
func listReqsById(id string) (*Requests, error) {
	log.Debugf("list requests %s", id)
	r := &Requests{}

	ops := ldb.Options{
		TableName:   "requests",
		Pks:         []string{"Id"},
		Ccms:        []string{},
		Hosts:       meta.MC.Scylla,
		Keyspace:    meta.MC.ScyllaKeyspace,
		PksClauses:  map[string]interface{}{"Id": id},
		CcmsClauses: make(map[string]interface{}),
	}
	if err := ldb.Fetchdb(ops, r); err != nil {
		return nil, err
	}

	log.Debugf("Requests %v", r)
	return r, nil
}
