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
package docker

import (
	log "code.google.com/p/log4go"
	"github.com/megamsys/libgo/amqp"
	"github.com/megamsys/megamd/app"
	"github.com/megamsys/megamd/provisioner"
	"encoding/json"
)

func Init() {
	provisioner.Register("docker", &Docker{})
}

type Message struct {
	Id string `json:"id"`
}

type Docker struct {
}

func (i *Docker) CreateCommand(assembly *provisioner.AssemblyResult, id string) (string, error) {
	predef := assembly.Components[0].Requirements.Host
	
	pdc, _ := app.GetPredefClouds(predef)
	assem := &app.Assembly{Id: pdc.Spec.Groups}
    dockerassembly, _ := assem.Get(pdc.Spec.Groups)	
	
    address := "Docker."+dockerassembly.Name+"."+dockerassembly.Components[0].Inputs.Domain
    com := &Message{Id: id}
	mapB, err := json.Marshal(com)  
	if err != nil {
        log.Error(err)
        return "", err
    }
    go publisher(address, string(mapB))
	return "", nil
}
func (i *Docker) DeleteCommand(assembly *provisioner.AssemblyResult, id string) (string, error) {
	
	return "", nil
}

func publisher(key string, json string) {
	factor, aerr := amqp.Factory()
	if aerr != nil {
		log.Error("Failed to get the queue instance: %s", aerr)
	}
	
	pubsub, perr := factor.Get(key)
	if perr != nil {
		log.Error("Failed to get the queue instance: %s", perr)
	}

	serr := pubsub.Pub([]byte(json))
	if serr != nil {
		log.Error("Failed to publish the queue instance: %s", serr)
	}
}