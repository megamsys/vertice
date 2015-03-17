package gogs

import (
  "github.com/megamsys/megamd/plugins"
  "github.com/megamsys/megamd/global"
  log "code.google.com/p/log4go"
  gogs "github.com/gogits/go-gogs-client"
  "strings"
  "encoding/json"
  "github.com/megamsys/libgo/amqp"
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
func (c *GogsPlugin) Watcher(ci *global.CI) error {
	if(ci.SCM == "gogs") {
		log.Info("gogs is worked")
		
		trigger_url := "https://api.megam.co/v2/assembly/build/"+ci.AssemblyID + "/" + ci.ComponentID 
		
		client := gogs.NewClient("http://192.168.1.5:6001/", ci.Token)
		
		var postData = make(map[string]string)
		postData["url"] = trigger_url
		postData["content_type"] = "json"
		
		postHook :=  gogs.CreateHookOption{Type: "gogs", Config: postData, Active: true }
		
		component := global.Component{Id: ci.ComponentID }
        com, comerr := component.Get(ci.ComponentID)
        if comerr != nil{
          return comerr
        }  
       
		source := strings.Split(com.Inputs.Source, "/") 
		
		s, _ := client.CreateRepoHook(ci.Owner, strings.Replace(source[len(source)-1], ".git", "", -1), postHook)
		
		//s, _ := client.ListRepoHooks(ci.Owner, strings.Replace(source[len(source)-1], ".git", "", -1))
		
		log.Info("---------------------------------------------")
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

