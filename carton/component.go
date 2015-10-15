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
	"time"

	"github.com/megamsys/megamd/db"
	"github.com/megamsys/megamd/provision"
	"github.com/megamsys/megamd/repository"
	"gopkg.in/yaml.v2"
)

const (
	DOMAIN        = "domain"
	PUBLICIP      = "publicip"
	BUCKET        = "components"
	IMAGE_VERSION = "version"
)

type Artifacts struct {
	Type         string    `json:"artifact_type"`
	Content      string    `json:"content"`
	Requirements JsonPairs `json:"requirements"`
}

/* Repository represents a repository managed by the manager. */
type Repo struct {
	Rtype    string `json:"rtype"`
	Source   string `json:"source"`
	Oneclick string `json:"oneclick"`
	Rurl     string `json:"url"`
}

type Component struct {
	Id                string        `json:"id"`
	Name              string        `json:"name"`
	Tosca             string        `json:"tosca_type"`
	Inputs            JsonPairs     `json:"inputs"`
	Outputs           JsonPairs     `json:"outputs"`
	Repo              Repo          `json:"repo"`
	Artifacts         *Artifacts    `json:"artifacts"`
	RelatedComponents []string      `json:"related_components"`
	Operations        []*Operations `json:"operations"`
	Status            string        `json:"status"`
	CreatedAt         string        `json:"created_at"`
}

func (a *Component) String() string {
	if d, err := yaml.Marshal(a); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

/**
**fetch the component json from riak and parse the json to struct
**/
func NewComponent(id string) (*Component, error) {
	c := &Component{Id: id}
	if err := db.Fetch(BUCKET, id, c); err != nil {
		return nil, err
	}
	return c, nil
}

//make a box with the details for a provisioner.
func (c *Component) mkBox() (provision.Box, error) {
	bt := provision.Box{
		Id:         c.Id,
		Level:      provision.BoxSome,
		Name:       c.Name,
		DomainName: c.domain(),
		Tosca:      c.Tosca,
		Commit:     "",
		Provider:   c.provider(),
		PublicIp:   c.publicIp(),
	}
	if &c.Repo != nil {
		bt.Repo = &repository.Repo{
			Type:     c.Repo.Rtype,
			Source:   c.Repo.Source,
			OneClick: c.Repo.Oneclick,
			URL:      c.Repo.Rurl,
		}
		bt.Repo.Hook = BuildHook(c.Operations, repository.CIHOOK)
	}
	return bt, nil
}

func (c *Component) SetStatus(status provision.Status) error {
	LastStatusUpdate := time.Now().In(time.UTC)
	c.Inputs = append(c.Inputs, NewJsonPair("lastsuccessstatusupdate", LastStatusUpdate.String()))
	c.Inputs = append(c.Inputs, NewJsonPair("status", status.String()))

	if err := db.Store(BUCKET, c.Id, c); err != nil {
		return err
	}
	return nil
}

func (c *Component) setDeployData(dd DeployData) error {
	c.Inputs = append(c.Inputs, NewJsonPair("lastsuccessstatusupdate", ""))
	c.Inputs = append(c.Inputs, NewJsonPair("status", ""))

	if err := db.Store(BUCKET, c.Id, c); err != nil {
		return err
	}
	return nil

}

func (c *Component) domain() string {
	return c.Inputs.match(DOMAIN)
}

func (c *Component) provider() string {
	return c.Inputs.match(provision.PROVIDER)
}

func (c *Component) publicIp() string {
	return c.Outputs.match(PUBLICIP)
}
