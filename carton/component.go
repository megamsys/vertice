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
	ldb "github.com/megamsys/libgo/db"
	"github.com/megamsys/libgo/pairs"
	"github.com/megamsys/libgo/utils"
	"github.com/megamsys/vertice/carton/bind"
	"github.com/megamsys/vertice/meta"
	"github.com/megamsys/vertice/provision"
	"github.com/megamsys/vertice/repository"
	"gopkg.in/yaml.v2"
	"strings"
	"time"
)

const (
	DOMAIN        = "domain"
	PUBLICIPV4    = "publicipv4"
	PRIVATEIPV4   = "privateipv4"
	COMPBUCKET    = "components"
	IMAGE_VERSION = "version"
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
	Source   string `json:"source" cql:"source"`
	Oneclick string `json:"oneclick" cql:"oneclick"`
	Rurl     string `json:"url" cql:"url"`
}

type Component struct {
	Id                string          `json:"id"`
	Name              string          `json:"name"`
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

type ComponentTable struct {
	Id                string   `json:"id" cql:"id"`
	Name              string   `json:"name" cql:"name"`
	Tosca             string   `json:"tosca_type" cql:"tosca_type"`
	Inputs            []string `json:"inputs" cql:"inputs"`
	Outputs           []string `json:"outputs" cql:"outputs"`
	Envs              []string `json:"envs" cql:"envs"`
	Repo              string   `json:"repo" cql:"repo"`
	Artifacts         string   `json:"artifacts" cql:"artifacts"`
	RelatedComponents []string `json:"related_components" cql:"related_components"`
	Operations        []string `json:"operations" cql:"operations"`
	Status            string   `json:"status" cql:"status"`
	State             string   `json:"state" cql:"state"`
	CreatedAt         string   `json:"created_at" cql:"created_at"`
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
	c := &ComponentTable{Id: id}
	ops := ldb.Options{
		TableName:   COMPBUCKET,
		Pks:         []string{"Id"},
		Ccms:        []string{},
		Hosts:       meta.MC.Scylla,
		Keyspace:    meta.MC.ScyllaKeyspace,
		PksClauses:  map[string]interface{}{"Id": id},
		CcmsClauses: make(map[string]interface{}),
	}
	if err := ldb.Fetchdb(ops, c); err != nil {
		return nil, err
	}
	com, _ := c.dig()
	return &com, nil
}

//make a box with the details for a provisioner.
func (c *Component) mkBox(vnet map[string]string,vmid string) (provision.Box, error) {
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
		Vnets:      vnet,
		VMId:       vmid,
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

func (c *Component) SetStatus(status utils.Status) error {
	LastStatusUpdate := time.Now().Local().Format(time.RFC822)
	m := make(map[string][]string, 2)
	m["lastsuccessstatusupdate"] = []string{LastStatusUpdate}
	m["status"] = []string{status.String()}
	c.Inputs.NukeAndSet(m) //just nuke the matching output key:

	update_fields := make(map[string]interface{})
	update_fields["Inputs"] = c.Inputs.ToString()
	update_fields["Status"] = status.String()
	ops := ldb.Options{
		TableName:   COMPBUCKET,
		Pks:         []string{"Id"},
		Ccms:        []string{},
		Hosts:       meta.MC.Scylla,
		Keyspace:    meta.MC.ScyllaKeyspace,
		PksClauses:  map[string]interface{}{"Id": c.Id},
		CcmsClauses: make(map[string]interface{}),
	}
	if err := ldb.Updatedb(ops, update_fields); err != nil {
		return err
	}
	return nil
}

func (c *Component) SetState(state utils.State) error {
	update_fields := make(map[string]interface{})
	update_fields["State"] = state.String()
	ops := ldb.Options{
		TableName:   COMPBUCKET,
		Pks:         []string{"Id"},
		Ccms:        []string{},
		Hosts:       meta.MC.Scylla,
		Keyspace:    meta.MC.ScyllaKeyspace,
		PksClauses:  map[string]interface{}{"Id": c.Id},
		CcmsClauses: make(map[string]interface{}),
	}
	if err := ldb.Updatedb(ops, update_fields); err != nil {
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
	ops := ldb.Options{
		TableName:   COMPBUCKET,
		Pks:         []string{"id"},
		Ccms:        []string{},
		Hosts:       meta.MC.Scylla,
		Keyspace:    meta.MC.ScyllaKeyspace,
		PksClauses:  map[string]interface{}{"id": compid},
		CcmsClauses: make(map[string]interface{}),
	}
	if err := ldb.Deletedb(ops, ComponentTable{}); err != nil {
		return
	}
}

func (c *Component) setDeployData(dd DeployData) error {
	/*c.Inputs = append(c.Inputs, utils.NewJsonPair("lastsuccessstatusupdate", ""))
	c.Inputs = append(c.Inputs, utils.NewJsonPair("status", ""))

	if err := db.Store(COMPBUCKET, c.Id, c); err != nil {
		return err
	}*/
	return nil

}

func (a *ComponentTable) dig() (Component, error) {
	asm := Component{}
	asm.Id = a.Id
	asm.Name = a.Name
	asm.Tosca = a.Tosca
	asm.Inputs = a.getInputs()
	asm.Outputs = a.getOutputs()
	asm.Envs = a.getEnvs()
	asm.Repo = a.getRepo()
	asm.Artifacts = a.getArtifacts()
	asm.RelatedComponents = a.RelatedComponents
	asm.Operations = a.getOperations()
	asm.Status = a.Status
	asm.CreatedAt = a.CreatedAt
	return asm, nil
}

func (c *Component) domain() string {
	return c.Inputs.Match(DOMAIN)
}

func (c *Component) provider() string {
	return c.Inputs.Match(utils.PROVIDER)
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

func (a *ComponentTable) getInputs() pairs.JsonPairs {
	keys := make([]*pairs.JsonPair, 0)
	for _, in := range a.Inputs {
		inputs := pairs.JsonPair{}
		parseStringToStruct(in, &inputs)
		keys = append(keys, &inputs)
	}
		return keys
}

func (a *ComponentTable) getOutputs() pairs.JsonPairs {
	keys := make([]*pairs.JsonPair, 0)
	for _, in := range a.Outputs {
		outputs := pairs.JsonPair{}
		parseStringToStruct(in, &outputs)
		keys = append(keys, &outputs)
	}
	return keys
}

func (a *ComponentTable) getEnvs() pairs.JsonPairs {
	keys := make([]*pairs.JsonPair, 0)
	for _, in := range a.Envs {
		outputs := pairs.JsonPair{}
		parseStringToStruct(in, &outputs)
		keys = append(keys, &outputs)
	}
	return keys
}

func (a *ComponentTable) getRepo() Repo {
	outputs := Repo{}
	parseStringToStruct(a.Repo, &outputs)
	return outputs
}

func (a *ComponentTable) getArtifacts() *Artifacts {
	outputs := Artifacts{}
	parseStringToStruct(a.Artifacts, &outputs)
	return &outputs
}

func (a *ComponentTable) getOperations() []*Operations {
	keys := make([]*Operations, 0)
	for _, in := range a.Operations {
		p := Operations{}
		parseStringToStruct(in, &p)
		keys = append(keys, &p)
	}
	return keys
}
