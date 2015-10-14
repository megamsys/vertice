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
	"github.com/megamsys/megamd/db"
	"github.com/megamsys/megamd/provision"
	"github.com/megamsys/megamd/repository"
	"gopkg.in/yaml.v2"
	"time"
)

const (
	DOMAIN        = "domain"
	PUBLICIP      = "publicip"
	BUCKET        = "components"
	IMAGE_VERSION = "version"
)

type Operations struct {
	Type        string    `json:"operation_type"`
	Description string    `json:"description"`
	Properties  JsonPairs `json:"properties"`
}

type Artifacts struct {
	Type         string    `json:"artifact_type"`
	Content      string    `json:"content"`
	Requirements JsonPairs `json:"requirements"`
}

type Component struct {
	Id                string           `json:"id"`
	Name              string           `json:"name"`
	Tosca             string           `json:"tosca_type"`
	Inputs            JsonPairs        `json:"inputs"`
	Outputs           JsonPairs        `json:"outputs"`
	Repo              *repository.Repo `json:"repo"`
	Artifacts         *Artifacts       `json:"artifacts"`
	RelatedComponents []string         `json:"related_components"`
	Operations        []*Operations    `json:"operations"`
	Status            string           `json:"status"`
	CreatedAt         string           `json:"created_at"`
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
	c.ciProps(c.Operations)

	return provision.Box{
		Id:         c.Id,
		Level:      provision.BoxSome,
		Name:       c.Name,
		DomainName: c.domain(),
		Tosca:      c.Tosca,
		Commit:     "",
		Repo:       c.Repo,
		Provider:   c.provider(),
		PublicIp:   c.publicIp(),
	}, nil

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

func (c *Component) ciProps(ops []*Operations) {
	o := parseOps(ops)

	if o != nil {
		c.Repo.Enabled = true
		c.Repo.Type = o.Properties.match(repository.TYPE)
		c.Repo.Token = o.Properties.match(repository.TOKEN)
		c.Repo.UserName = o.Properties.match(repository.USERNAME)
	}
}

func parseOps(ops []*Operations) *Operations {
	for _, o := range ops {
		switch o.Type {
		case repository.OPS_CI:
			return o
		}
	}
	return nil
}
