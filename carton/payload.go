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
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/megamd/db"
)

type Payload struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CatId     string `json:"cat_id"`
	CatType   string `json:"cattype"`
	CreatedAt string `json:"created_at"`
}

type PayloadConvertor interface {
	Convert(p *Payload) (*Requests, error)
}

func NewPayload(b []byte) (*Payload, error) {
	p := &Payload{}
	err := json.Unmarshal(b, &p)
	if err != nil {
		log.Error("Failed to parse the payload message:\n%s.", err)
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
			Name:      p.Name,
			CatId:     p.CatId,
			CatType:   p.CatType,
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
	if err := db.Fetch("requests", id, r); err != nil {
		return nil, err
	}

	log.Debugf("Requests %v", r)
	return r, nil
}
