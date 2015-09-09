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
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/megamd/db"
)

type Payload struct {
	Id string `json:"id"`
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
	log.Infof("get requests %s", p.Id)
	r := &Requests{}
	if err := db.Fetch("requests", p.Id, r); err != nil {
		return nil, err
	}

	log.Debugf("Requests %v", r)
	return r, nil

}
