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
	"strconv"
	"strings"
	"time"
)

const (
	DOMAIN = "domain"
	BUCKET = "components"
)

type Operations struct {
	OperationType         string    `json:"operation_type"`
	Description           string    `json:"description"`
	OperationRequirements JsonPairs `json:"operation_requirements"`
}

type Artifacts struct {
	ArtifactType         string    `json:"artifact_type"`
	Content              string    `json:"content"`
	ArtifactRequirements JsonPairs `json:"artifact_requirements"`
}

type Component struct {
	Id                string        `json:"id"`
	Name              string        `json:"name"`
	Tosca             string        `json:"tosca_type"`
	Inputs            JsonPairs     `json:"inputs"`
	Outputs           JsonPairs     `json:"outputs"`
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
	repo := c.NewRepo(repository.CI)

	return provision.Box{
		Id:         c.Id,
		Level:      provision.BoxSome,
		Name:       c.Name,
		DomainName: c.domain(),
		Tosca:      c.Tosca,
		Commit:     "",
		Image:      c.image(),
		Repo:       repo,
		Provider:   c.provider(),
		Ip:         "",
	}, nil
}

func (c *Component) SetStatus(status provision.Status) error {
	LastStatusUpdate := time.Now().In(time.UTC)

	/*
	   //masking this check for now. why do we need this status check ?
	   if c.Status == provision.StatusDeploying.String() ||
	   		c.Status == provision.StatusCreating.String() ||
	   		c.Status == provision.StatusCreated.String() ||
	   		c.Status == provision.StatusStateup.String() {
	*/
	c.Inputs = append(c.Inputs, NewJsonPair("lastsuccessstatusupdate", LastStatusUpdate.String()))
	c.Inputs = append(c.Inputs, NewJsonPair("status", status.String()))
	//	}

	if err := db.Store(BUCKET, c.Id, c); err != nil {
		return err
	}
	return nil

}

func (c *Component) NewRepo(ci string) repository.Repo {
	o := parseOps(c.Operations, ci)

	if o != nil {
		enabled, _ := strconv.ParseBool(o.OperationRequirements.match(repository.CI_ENABLED))

		return repository.Repo{
			Enabled:     enabled,
			Token:       o.OperationRequirements.match(repository.CI_TOKEN),
			Git:         o.OperationRequirements.match(repository.CI_SCM),
			UserName:    o.OperationRequirements.match(repository.CI_USER),			
		}

	}
	return repository.Repo{}

}

func (c *Component) domain() string {
	return c.Inputs.match(DOMAIN)
}

func (c *Component) provider() string {
	return c.Inputs.match(provision.PROVIDER)
}

// for a vm provisioner return the last name (tosca.torpedo.ubuntu) ubuntu as the image name.
// for docker return the Inputs[image]
func (c *Component) image() string {
	switch c.provider() {
	case provision.PROVIDER_ONE:
		return c.Tosca[strings.LastIndex(c.Tosca, ".")+1:]
	case provision.PROVIDER_DOCKER:
		return c.Inputs.match("image")
	default:
		return ""
	}
}

func parseOps(ops []*Operations, optype string) *Operations {
	for _, o := range ops {
		switch o.OperationType {
		case repository.CI:
			return o
		}
	}
	return nil
}
