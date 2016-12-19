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
	"github.com/megamsys/libgo/api"
	"io/ioutil"
	"strings"
	"time"
)

type Payload struct {
	Id        string `json:"id"`
	Action    string `json:"action"`
	CatId     string `json:"cat_id"`
	AccountId string `json:"email"`
	CatType   string `json:"cattype"`
	Category  string `json:"category"`
	CreatedAt time.Time `json:"created_at"`
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
	if len(strings.TrimSpace(p.CatId)) < 1 {
		return listReqsById(p.Id, p.AccountId)
	} else {
		return &Requests{
			Action:    p.Action,
			Category:  p.Category,
			AccountId: p.AccountId,
			CatId:     p.CatId,
			CreatedAt: p.CreatedAt,
		}, nil
	}

}

//The payload in the queue can be just a pointer or a value.
//pointer means just the id will be available and rest is blank.
//value means the id is blank and others are available.
func listReqsById(id, email string) (*Requests, error) {
	log.Debugf("list requests %s", id)
	args := newArgs(email,"")
	args.Path = "/requests/" + id
	cl := api.NewClient(args)
	response, err := cl.Get()
	if err != nil {
		return nil, err
	}
	htmlData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	res := &ApiRequests{}
	err = json.Unmarshal(htmlData, res)
	if err != nil {
		return nil, err
	}
	r := &res.Results[0]
	log.Debugf("Requests %v", r)
	return r, nil
}
