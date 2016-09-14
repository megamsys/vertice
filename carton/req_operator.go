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
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
)

type ReqOperator struct {
	Id        string
	CartonsId string
	Category  string
	Action    string
}

// NewReqOperator returns a new instance of ReqOperator
// for the operatable id (Assemblies)
func NewReqOperator(id ,cat,action string) *ReqOperator {
	return &ReqOperator{Id: id, Category: cat, Action: action}
}

func (p *ReqOperator) Accept(r *MegdProcessor) error {
	c, err := p.Get()
	if err != nil {
		return err
	}
	md := *r
	log.Debugf(cmd.Colorfy(md.String(), "cyan", "", "bold"))
	return md.Process(c)
}

func (p *ReqOperator) Get() (Cartons, error) {
  if p.Category == SNAPSHOT {
		s, err := GetSnap(p.Id)
		if err != nil {
			return nil, err
		}
		c, err := s.MkCartons()
		if err != nil {
			return nil, err
		}
		return c, nil
	} else if p.Category == DISKS {
		d, err := GetDisks(p.Id)
		if err != nil {
			return nil, err
		}
		c, err := d.MkCartons()
		if err != nil {
			return nil, err
		}
		return c, nil
	} else {
		a, err := Get(p.Id)
		if err != nil {
			return nil, err
		}
		c, err := a.MkCartons()
		if err != nil {
			return nil, err
		}
		return c, nil
	}
}


// MegdProcessor represents a single operation in vertice.
type MegdProcessor interface {
	Process(c Cartons) error
	String() string
}
