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
	"github.com/megamsys/libgo/api"
	"github.com/megamsys/libgo/pairs"
	"github.com/megamsys/libgo/utils"
	"github.com/megamsys/vertice/carton/bind"
	"github.com/megamsys/vertice/provision"
	"github.com/megamsys/vertice/repository"
	"gopkg.in/yaml.v2"
	"strings"
	"time"
)

const (
	DOMAIN        = "domain"
	COMPBUCKET    = "components"
	IMAGE_VERSION = "version"
	BACKUPNAME    = "backup_name"
	ONECLICK      = "oneclick"
	HOSTIP        = "vnchost"
	VERTICE       = "vertice"
	TRUE          = "true"
)

type Artifacts struct {
	Type         string          `json:"artifact_type" cql:"type"`
	Content      string          `json:"content" cql:"content"`
	Requirements pairs.JsonPairs `json:"requirements" cql:"requirements"`
}

/* Repository represents a repository managed by the manager. */
type Repo struct {
	Rtype    string `json:"rtype" cql:"rtype"`
	Branch   string `json:"branch" cql:"branch"`
	Source   string `json:"source" cql:"source"`
	Oneclick string `json:"oneclick" cql:"oneclick"`
	Rurl     string `json:"url" cql:"url"`
}

type ApiComponent struct {
	JsonClaz string      `json:"json_claz"`
	Results  []Component `json:"results"`
}

type Component struct {
	Id                string          `json:"id"`
	Name              string          `json:"name"`
	OrgId             string          `json:"org_id"`
	Tosca             string          `json:"tosca_type"`
	Inputs            pairs.JsonPairs `json:"inputs"`
	Outputs           pairs.JsonPairs `json:"outputs"`
	Envs              pairs.JsonPairs `json:"envs"`
	Repo              Repo            `json:"repo"`
	Artifacts         *Artifacts      `json:"artifacts"`
	RelatedComponents []string        `json:"related_components"`
	Operations        []*Operations   `json:"operations"`
	Status            string          `json:"status"`
	State             string          `json:"state"`
	CreatedAt         string          `json:"created_at"`
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

func NewComponent(id, email, org string) (*Component, error) {
	cl := api.NewClient(newArgs(email, org), "/components/"+id)
	response, err := cl.Get()
	if err != nil {
		return nil, err
	}

	ac := &ApiComponent{}
	err = json.Unmarshal(response, ac)
	if err != nil {
		return nil, err
	}
	return &ac.Results[0], nil
}

func (c *Component) updateComponent(email, org string) error {
	cl := api.NewClient(newArgs(email, org), "/components/update")
	_, err := cl.Post(c)
	if err != nil {
		return err
	}
	return nil
}

//make a box with the details for a provisioner.
func (c *Component) mkBox() (provision.Box, error) {
	bt := provision.Box{
		Id:          c.Id,
		Level:       provision.BoxSome,
		Name:        c.Name,
		DomainName:  c.domain(),
		Envs:        c.envs(),
		Tosca:       c.Tosca,
		Commit:      "",
		Provider:    c.provider(),
		PublicIp:    c.publicIp(),
		StorageType: c.storageType(),
		OrgId:       c.OrgId,
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

func (c *Component) SetStatus(status utils.Status, email string) error {
	LastStatusUpdate := time.Now().Local().Format(time.RFC822)
	m := make(map[string][]string, 2)
	m["lastsuccessstatusupdate"] = []string{LastStatusUpdate}
	m["status"] = []string{status.String()}
	c.Inputs.NukeAndSet(m) //just nuke the matching output key:
	c.Status = status.String()
	return c.updateComponent(email, c.OrgId)
}

func (c *Component) SetState(state utils.State, email string) error {
	c.State = state.String()
	return c.updateComponent(email, c.OrgId)
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

func (c *Component) Delete(email, orgid string) error {
	cl := api.NewClient(newArgs(email, orgid), "/components/"+c.Id)
	_, err := cl.Delete()
	if err != nil {
		return err
	}
	return nil
}

func (c *Component) setDeployData(dd DeployData) error {
	/*c.Inputs = append(c.Inputs, utils.NewJsonPair("lastsuccessstatusupdate", ""))
	c.Inputs = append(c.Inputs, utils.NewJsonPair("status", ""))

	if err := db.Store(COMPBUCKET, c.Id, c); err != nil {
		return err
	}*/
	return nil

}

func (c *Component) domain() string {
	return c.Inputs.Match(DOMAIN)
}

func (c *Component) provider() string {
	return c.Inputs.Match(utils.PROVIDER)
}

func (c *Component) storageType() string {
	return strings.ToLower(c.Inputs.Match(utils.STORAGE_TYPE))
}

func (c *Component) publicIp() string {
	return c.Outputs.Match(PUBLICIPV4)
}

func (c *Component) withOneClick() bool {
	return (strings.TrimSpace(c.Repo.Oneclick) == TRUE && c.Repo.Source == VERTICE)
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
