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
package coordinator

import (
	log "code.google.com/p/log4go"
	"encoding/json"
	"github.com/megamsys/megamd/app"
	"github.com/megamsys/megamd/iaas/ec2"
	"github.com/megamsys/megamd/iaas/gogrid"
	"github.com/megamsys/megamd/iaas/google"
	"github.com/megamsys/megamd/iaas/hp"
	"github.com/megamsys/megamd/iaas/opennebula"
	"github.com/megamsys/megamd/iaas/profitbricks"
	"github.com/megamsys/megamd/provisioner/chef"
	"github.com/megamsys/megamd/provisioner/docker"
	"github.com/megamsys/megamd/plugins/cmp"
	"github.com/megamsys/megamd/plugins/github"
	"github.com/megamsys/megamd/plugins/gogs"
	"github.com/megamsys/megamd/iaas/megam"
	"github.com/megamsys/megamd/plugins"
	"github.com/megamsys/megamd/global"
	"fmt"
)

type Coordinator struct {
	//RequestHandler(f func(*Message), name ...string) (Handler, error)
	//EventsHandler(f func(*Message), name ...string) (Handler, error)
}

func init() {
	ec2.Init()
	google.Init()
	hp.Init()
	gogrid.Init()
	opennebula.Init()
	profitbricks.Init()
	chef.Init()
	docker.Init()
	cmp.Init()
	github.Init()
	gogs.Init()
	megam.Init()
}

func NewCoordinator(chann []byte, queue string) {
	log.Info("Handling coordinator message %v", string(chann))

	switch queue {
	case "cloudstandup":
		requestHandler(chann)
		break
	case "events":
		eventsHandler(chann)
		break
	case "ci":
	    ciHandler(chann)
	    break	
	}
}

func requestHandler(chann []byte) {
	log.Info("Entered!-------->")
	m := &global.Message{}
	parse_err := json.Unmarshal(chann, &m)
	if parse_err != nil {
		log.Error("Error: Message parsing error:\n%s.", parse_err)
		return
	}
	request := app.Request{Id: m.Id}
	req, err := request.Get(m.Id)
	log.Debug(req)
	log.Debug("---------")
	if err != nil {
		log.Error("Error: Riak didn't cooperate:\n%s.", err)
		return
	}
	switch req.ReqType {
	case "create":
	log.Debug("============Create entry======")
		assemblies := app.Assemblies{Id: req.AssembliesId}
		asm, err := assemblies.Get(req.AssembliesId)
		log.Debug("----------")
		log.Debug(asm)
		log.Debug("------------")
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return
		}
		for i := range asm.Assemblies {
			log.Debug("Assemblies: [%s]", asm.Assemblies[i])
			if len(asm.Assemblies[i]) > 1 {
				assemblyID := asm.Assemblies[i]
				log.Debug("Assemblies id: [%s]", assemblyID)
				assembly := app.Assembly{Id: assemblyID}
				log.Debug(assembly)
				res, err := assembly.Get(assemblyID)
				log.Debug(res)
				if err != nil {
					log.Error("Error: Riak didn't cooperate:\n%s.", err)
					return
				}              
				go app.LaunchApp(res, m.Id)
			}
		}
		
		//build delete command
	    case "delete":
		log.Debug("============Delete entry==========")
		  assembly := app.Assembly{Id: req.AssembliesId}
		  asm, err := assembly.Get(req.AssembliesId)
		   log.Debug("----------")
		   log.Debug(asm)
		   log.Debug("------------")
		   if err != nil {
		   	   log.Error("Error: Riak didn't cooperate:\n%s.", err)
		   	   return
		   }
		   res := asm
		   	   	 go app.DeleteApp(res, m.Id)
		   	   
		   	   
		   
	}
}

func ciHandler(chann []byte) {
	log.Info("Entered!-------->")
	m := &global.Message{}
	parse_err := json.Unmarshal(chann, &m)
	if parse_err != nil {
		log.Error("Error: Message parsing error:\n%s.", parse_err)
		return
	}
	request := global.CI{Id: m.Id}
	cig, err := request.Get(m.Id)
	fmt.Println("-----====================================----")
	fmt.Println(cig)
	
	perr := plugins.Watcher(cig)
	if perr != nil {
		log.Error("Error: Plugin Watcher :\n%s.", perr)
		return
	}
	
	if err != nil {
		log.Error("Error: Riak didn't cooperate:\n%s.", err)
		return
	}
}

func eventsHandler(chann []byte) {
   log.Info("Event was entered")
   m := &global.EventMessage{}
	parse_err := json.Unmarshal(chann, &m)
	if parse_err != nil {
		log.Error("Error: Message parsing error:\n%s.", parse_err)
		return
	}
	switch m.Event {
	case "notify":
		perr := plugins.Notify(m)
		if perr != nil {
			log.Error("Error: Plugin Notify :\n%s.", perr)
			break
		}
		break	
	}
}
