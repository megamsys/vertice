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
	"time"

	"github.com/megamsys/megamd/carton/bind"
	"github.com/megamsys/megamd/db"
	"github.com/megamsys/megamd/provision"
	"github.com/megamsys/megamd/repository"
	"gopkg.in/yaml.v2"
)

const (
	DOMAIN        = "domain"
	PUBLICIPV4    = "publicipv4"
	PRIVATEIPV4   = "privateipv4"
	COMPBUCKET    = "components"
	IMAGE_VERSION = "version"
	ONECLICK      = "oneclick"
)

type Artifacts struct {
	Type         string         `json:"artifact_type"`
	Content      string         `json:"content"`
	Requirements bind.JsonPairs `json:"requirements"`
}

/* Repository represents a repository managed by the manager. */
type Repo struct {
	Rtype    string `json:"rtype"`
	Source   string `json:"source"`
	Oneclick string `json:"oneclick"`
	Rurl     string `json:"url"`
}

type Component struct {
	Id                string         `json:"id"`
	Name              string         `json:"name"`
	Tosca             string         `json:"tosca_type"`
	Inputs            bind.JsonPairs `json:"inputs"`
	Outputs           bind.JsonPairs `json:"outputs"`
	Envs              bind.JsonPairs `json:"envs"`
	Repo              Repo           `json:"repo"`
	Artifacts         *Artifacts     `json:"artifacts"`
	RelatedComponents []string       `json:"related_components"`
	Operations        []*Operations  `json:"operations"`
	Status            string         `json:"status"`
	CreatedAt         string         `json:"created_at"`
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
	if err := db.Fetch(COMPBUCKET, id, c); err != nil {
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
		Envs:       c.envs(),
		Tosca:      c.Tosca,
		Commit:     "",
		Provider:   c.provider(),
		PublicIp:   c.publicIp(),
	}

	if &c.Repo != nil {
		bt.Repo = &repository.Repo{
			Type:     c.Repo.Rtype,
			Source:   c.Repo.Source,
			OneClick: c.withOneClick(),
			URL:      c.Repo.Rurl,
		}
		bt.Repo.Hook = BuildHook(c.Operations, repository.CIHOOK)
	}
	return bt, nil
}

func (c *Component) SetStatus(status provision.Status) error {
	LastStatusUpdate := time.Now().Local().Format(time.RFC822)
	m := make(map[string][]string, 2)
	m["lastsuccessstatusupdate"] = []string{LastStatusUpdate}
	m["status"] = []string{status.String()}
	c.Inputs.NukeAndSet(m) //just nuke the matching output key:

	c.Status = status.String()

	if err := db.Store(COMPBUCKET, c.Id, c); err != nil {
		return err
	}
	return nil
}

/*func (c *Component) UpdateOpsRun(opsRan upgrade.OperationsRan) error {
	mutatedOps := make([]*upgrade.Operation, 0, len(opsRan))

	for _, o := range opsRan {
		mutatedOps = append(mutatedOps, o.Raw)
	}
	c.Operations = mutatedOps

	if err := db.Store(COMPBUCKET, c.Id, c); err != nil {
		return err
	}
	return nil
}*/

func (c *Component) Delete(compid string) {
	_ = db.Delete(COMPBUCKET, compid)
}

func (c *Component) setDeployData(dd DeployData) error {
	/*c.Inputs = append(c.Inputs, bind.NewJsonPair("lastsuccessstatusupdate", ""))
	c.Inputs = append(c.Inputs, bind.NewJsonPair("status", ""))

	if err := db.Store(COMPBUCKET, c.Id, c); err != nil {
		return err
	}*/
	return nil

}

func (c *Component) domain() string {
	return c.Inputs.Match(DOMAIN)
}

func (c *Component) provider() string {
	return c.Inputs.Match(provision.PROVIDER)
}

func (c *Component) publicIp() string {
	return c.Outputs.Match(PUBLICIPV4)
}

func (c *Component) withOneClick() bool {
	return (len(strings.TrimSpace(c.Envs.Match(ONECLICK))) > 0)
}

//all the variables in the inputs shall be treated as ENV.
//we can use a filtered approach as well.
func (c *Component) envs() []bind.EnvVar {
	envs := make([]bind.EnvVar, 0, len(c.Envs))
	for _, i := range c.Envs {
		envs = append(envs, bind.EnvVar{Name: i.K, Value: i.V})
	}
	return envs
}
