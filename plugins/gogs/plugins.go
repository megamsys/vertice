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
package gogs

import (
  "github.com/megamsys/megamd/plugins"
  "github.com/megamsys/megamd/global"
  log "code.google.com/p/log4go"
  gogs "github.com/gogits/go-gogs-client"
  "strings"
  "github.com/tsuru/config"
)

const ( 
   GOGS = "gogs"
   ENABLE = "true"   
)

func Init() {
	plugins.RegisterPlugins("gogs", &GogsPlugin{})
}

type GogsPlugin struct{}

/**
**watcher function executes the gogs repository webhook creation
**first get the ci value and parse it and to create the hook for that users repository    
**get trigger url from config file 
**/
func (c *GogsPlugin) Watcher(asm *global.AssemblyWithComponents, ci *global.Operations, com *global.Component) error {
     switch ci.OperationType {
		case "CI":
			cierr := cioperation(asm, ci, com)
			if cierr != nil {
					return cierr
			}
			break
     }     		
	return nil
}

func cioperation(asm *global.AssemblyWithComponents, ci *global.Operations, com *global.Component) error {
	 pair_scm, perrscm := global.ParseKeyValuePair(ci.OperationRequirements, "ci-scm")
		if perrscm != nil {
			log.Error("Failed to get the domain value : %s", perrscm)
		}
		
	pair_enable, perrenable := global.ParseKeyValuePair(ci.OperationRequirements, "ci-enable")
		if perrenable != nil {
			log.Error("Failed to get the domain value : %s", perrenable)
		}	
		
	if(pair_scm.Value == GOGS && pair_enable.Value == ENABLE) {
		log.Info("gogs process started...")
		
		//trigger_url := "https://api.megam.co/v2/assembly/build/"+asm.Id + "/" + com.Id
		trigger_url := "http://localhost:9000/v2/assembly/build/"+asm.Id + "/" + com.Id
		
		url, herr := config.GetString("gogs:url")		
		if herr != nil {		  
			return herr
		}		
		
		pair_token, perrtoken := global.ParseKeyValuePair(ci.OperationRequirements, "ci-token")
		if perrtoken != nil {
			log.Error("Failed to get the domain value : %s", perrtoken)
		}
		
		client := gogs.NewClient(url, pair_token.Value)
		log.Info("Gogs api client created")
		
		var postData = make(map[string]string)
		postData["url"] = trigger_url
		postData["content_type"] = "json"
		
		postHook :=  gogs.CreateHookOption{Type: "gogs", Config: postData, Active: true }
		
		pair_source, perrsource := global.ParseKeyValuePair(com.Inputs, "source")
		if perrsource != nil {
			log.Error("Failed to get the domain value : %s", perrsource)
		}
		       
		source := strings.Split(pair_source.Value, "/") 
		log.Info(strings.Replace(source[len(source)-1], ".git", "", -1))
		
		pair_owner, perrowner := global.ParseKeyValuePair(ci.OperationRequirements, "ci-owner")
		if perrowner != nil {
			log.Error("Failed to get the domain value : %s", perrowner)
		}
		
		s, hook_err := client.CreateRepoHook(pair_owner.Value, strings.Replace(source[len(source)-1], ".git", "", -1), postHook)
		if hook_err !=nil {		
		   return hook_err
		}
		//s, _ := client.ListRepoHooks(ci.Owner, strings.Replace(source[len(source)-1], ".git", "", -1))
		
		log.Info("Hook created")
		log.Info(s)		
	} else {
		log.Info("gogs is skipped")
	}
	
	return nil
}

/**
**notify the messages or any other operations from github
** now this function performs build the pushed application from gogs to remote 
**/
func (c *GogsPlugin) Notify(m *global.EventMessage) error {
	/*request_asm := global.Assembly{Id: m.AssemblyId}
	asm, asmerr := request_asm.Get(m.AssemblyId)
	if(asmerr != nil) {
		return asmerr
	}

	request_com := global.Component{Id: m.ComponentId}
	com, comerr := request_com.Get(m.ComponentId)
	if(comerr != nil) {
		return comerr
	}
	
	request_ci := global.CI{Id: com.Inputs.CIID}
	ci, cierr := request_ci.Get(com.Inputs.CIID)
	if(cierr != nil) {
		return cierr
	}
	
	if(ci.SCM == "gogs") {
		log.Info("Gogs is working")
		mapD := map[string]string{"Id": m.ComponentId, "Action": "build"}
		mapB, _ := json.Marshal(mapD)
		log.Info(string(mapB))
		asmname := asm.Name
		//asmname := asm.Name
		publisher(asmname, string(mapB))
	} else {
		log.Info("Gogs is skipped")
	}*/
	return nil
}

func publisher(key string, json string) {
	/*factor, aerr := amqp.Factory()
	if aerr != nil {
		log.Error("Failed to get the queue instance: %s", aerr)
	}
	//s := strings.Split(key, "/")
	//pubsub, perr := factor.Get(s[len(s)-1])
	pubsub, perr := factor.Get(key)
	if perr != nil {
		log.Error("Failed to get the queue instance: %s", perr)
	}

	serr := pubsub.Pub([]byte(json))
	if serr != nil {
		log.Error("Failed to publish the queue instance: %s", serr)
	}*/
}

