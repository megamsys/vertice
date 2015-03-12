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
		
		trigger_url := "https://api.megam.co/v2/assembly/build/"+ci.AssemblyID 
		
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

func (c *GithubPlugin) Notify() error {
	return nil
}