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
package github

import (
  "github.com/megamsys/megamd/plugins"
  "github.com/megamsys/megamd/global"
  log "code.google.com/p/log4go"
  //"fmt"
  //git "github.com/CodeHub-io/Go-GitHub-API"
  git "github.com/google/go-github/github"  
  "code.google.com/p/goauth2/oauth"
  "encoding/json"
  "strings"
  "github.com/megamsys/libgo/amqp"
)

/**
**Github register function
**This function register to plugins container   
**/
func Init() {
	plugins.RegisterPlugins("github", &GithubPlugin{})
}

type GithubPlugin struct{}

/**
**watcher function executes the github repository webhook creation
**first get the ci value and parse it and to create the hook for that users repository    
**get trigger url from config file 
**/
func (c *GithubPlugin) Watcher(ci *global.CI) error {
	if(ci.SCM == "github") {
		log.Info("Github is working")
		
		t := &oauth.Transport{
			Token: &oauth.Token{AccessToken: ci.Token},
		}
		
		client := git.NewClient(t.Client())
		
		trigger_url := "https://api.megam.co/v2/assembly/build/"+ci.AssemblyID + "/" + ci.ComponentID 
		
		byt12 := []byte(`{"url": "","content_type": "json"}`)
		var postData map[string]interface{}
    	if perr := json.Unmarshal(byt12, &postData); perr != nil {
        	return perr
    	}
		postData["url"] = trigger_url
		
		byt1 := []byte(`{"name": "web", "active": true, "events": [ "push" ]}`)
		postHook :=  git.Hook{Config: postData }
    	if perr := json.Unmarshal(byt1, &postHook); perr != nil {
        	log.Info(perr)
    	}		
    	//postHook.Name = postHook.String(global.RandString(6))
    	component := global.Component{Id: ci.ComponentID }
        com, comerr := component.Get(ci.ComponentID)
        if comerr != nil{
          return comerr
        }  
        
        source := strings.Split(com.Inputs.Source, "/")    			
		_, _, err := client.Repositories.CreateHook(ci.Owner, strings.TrimRight(source[len(source)-1], ".git"), &postHook)
		
    	if err != nil {
        	return err
    	}
    	
	} else {
		log.Info("Github is skipped")
	}
	
	return nil
}

/**
**notify the messages or any other operations from github
** now this function performs build the pushed application from github to remote 
**/
func (c *GithubPlugin) Notify(m *global.EventMessage) error {
   request_asm := global.Assembly{Id: m.AssemblyId}
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
	
	if(ci.SCM == "github") {
		log.Info("Github is worked")
		mapD := map[string]string{"Id": m.ComponentId, "Action": "build"}
		mapB, _ := json.Marshal(mapD)
		log.Info(string(mapB))
		asmname := asm.Name
		//asmname := asm.Name
		publisher(asmname, string(mapB))
	} else {
		log.Info("Github is skipped")
	}
	return nil
}

func publisher(key string, json string) {
	factor, aerr := amqp.Factory()
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
	}
}
