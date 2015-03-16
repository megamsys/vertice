package gogs

import (
  "github.com/megamsys/megamd/plugins"
  "github.com/megamsys/megamd/global"
  log "code.google.com/p/log4go"
  gogs "github.com/gogits/go-gogs-client"
  "strings"
)


func Init() {
	plugins.RegisterPlugins("gogs", &GogsPlugin{})
}

type GogsPlugin struct{}


func (c *GogsPlugin) Watcher(ci *global.CI) error {
	if(ci.SCM == "gogs") {
		log.Info("gogs is worked")
		
		trigger_url := "https://api.megam.co/v2/assembly/build/"+ci.AssemblyID 
		
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

func (c *GogsPlugin) Notify() error {
	return nil
}