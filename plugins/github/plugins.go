package github

import (
  "github.com/megamsys/megamd/plugins"
  "github.com/megamsys/megamd/global"
  log "code.google.com/p/log4go"
)


func Init() {
	plugins.RegisterPlugins("github", &GithubPlugin{})
}

type GithubPlugin struct{}


func (c *GithubPlugin) Watcher(ci *global.CI) error {
	if(ci.SCM == "github") {
		log.Info("Github is worked")
	} else {
		log.Info("Github is skipped")
	}
	
	return nil
}

func (c *GithubPlugin) Notify() error {
	return nil
}