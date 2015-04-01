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
package provisioner

import (
	//	"errors"
	log "code.google.com/p/log4go"
	"fmt"
	"github.com/megamsys/libgo/db"
)

type Policy struct {
	Name    string   `json:"name"`
	Ptype   string   `json:"ptype"`
	Members []string `json:"members"`
}

type AssemblyResult struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Components []*Component
	policies   []*Policy `json:"policies"`
	inputs     string    `json:"inputs"`
	operations string    `json:"operations"`
	Command    string
	CreatedAt  string `json:"created_at"`
}

type Component struct {
	Id                         string `json:"id"`
	Name                       string `json:"name"`
	ToscaType                  string `json:"tosca_type"`
	Requirements               *ComponentRequirements
	Inputs                     *CompomentInputs
	ExternalManagementResource string
	Artifacts                  *Artifacts
	RelatedComponents          string
	Operations                 *ComponentOperations
	CreatedAt                  string `json:"created_at"`
}

type ComponentRequirements struct {
	Host  string `json:"host"`
	Dummy string `json:"dummy"`
}

type CompomentInputs struct {
	Domain        string `json:"domain"`
	Port          string `json:"port"`
	UserName      string `json:"username"`
	Password      string `json:"password"`
	Version       string `json:"version"`
	Source        string `json:"source"`
	DesignInputs  *DesignInputs
	ServiceInputs *ServiceInputs
	CIID          string  `json:"ci_id"`
}


type DesignInputs struct {
	Id    string   `json:"id"`
	X     string   `json:“x”`
	Y     string   `json:“y”`
	Z     string   `json:“z”`
	Wires []string `json:“wires”`
}

type ServiceInputs struct {
	DBName     string `json:"dbname"`
	DBPassword string `json:“dbpassword”`
}

type Artifacts struct {
	ArtifactType string `json:"artifact_type"`
	Content      string `json:“content”`
}

type ComponentOperations struct {
	OperationType  string `json:"operation_type"`
	TargetResource string `json:“target_resource”`
}

func (com *Component) Get(comId string) error {
	log.Info("Get message %v", comId)
	conn, err := db.Conn("components")
	if err != nil {
		return err
	}
	//appout := &Requests{}
	ferr := conn.FetchStruct(comId, com)
	if ferr != nil {
		return ferr
	}
	defer conn.Close()
	return nil

}

// Provisioner is the basic interface of this package.
//
// Any tsuru provisioner must implement this interface in order to provision
// tsuru apps.
type Provisioner interface {
	// Provision is called when tsuru is creating the app.
	//Provision(*AssemblyResult) error

	// ExecuteCommand runs a command in all units of the app.
	 CreateCommand(*AssemblyResult, string) (string, error)
     DeleteCommand(*AssemblyResult, string) (string, error)
	// ExecuteCommandOnce runs a command in one unit of the app.
	//	ExecuteCommandOnce(stdout, stderr io.Writer, app app.AssemblyResult, cmd string, args ...string) error

}

var provisioners = make(map[string]Provisioner)

// Register registers a new provisioner in the Provisioner registry.
func Register(name string, p Provisioner) {
	provisioners[name] = p
}

func GetProvisioner(name string) (Provisioner, error) {

	provider, ok := provisioners[name]
	if !ok {
		return nil, fmt.Errorf("Provisioner not registered")
	}
	return provider, nil
	//return nil
}
