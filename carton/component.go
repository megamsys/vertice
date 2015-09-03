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
	log "github.com/golang/glog"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/db"
)

type Operations struct {
	OperationType         string          `json:"operation_type"`
	Description           string          `json:"description"`
	OperationRequirements []*KeyValuePair `json:"operation_requirements"`
}

type Artifacts struct {
	ArtifactType         string          `json:"artifact_type"`
	Content              string          `json:"content"`
	ArtifactRequirements []*KeyValuePair `json:"artifact_requirements"`
}

type Component struct {
	Id                string          `json:"id"`
	Name              string          `json:"name"`
	ToscaType         string          `json:"tosca_type"`
	Inputs            []*KeyValuePair `json:"inputs"`
	Outputs           []*KeyValuePair `json:"outputs"`
	Artifacts         *Artifacts      `json:"artifacts"`
	RelatedComponents []string        `json:"related_components"`
	Operations        []*Operations   `json:"operations"`
	Status            string          `json:"status"`
	CreatedAt         string          `json:"created_at"`
}

func NewComponent(id string) *Component {
	&Component{Id: id}
}

/**
**fetch the component json from riak and parse the json to struct
**/
func (c *Component) Get(comp_id string) error {
	log.Debugf("[global] Get component %s", asmId)
	if conn, err := db.Conn("components"); err != nil {
		return err
	}

	if err := conn.FetchStruct(comp_id, c); err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func (c *Component) mkBox() (*provisioner.Box, error) {
	return nil, nil
}
