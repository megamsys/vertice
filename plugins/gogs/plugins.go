package gogs

import (
  "github.com/megamsys/megamd/plugins"
  "github.com/megamsys/megamd/global"
  log "code.google.com/p/log4go"
)


func Init() {
	plugins.RegisterPlugins("gogs", &GogsPlugin{})
}

type GogsPlugin struct{}


func (c *GogsPlugin) Watcher(ci *global.CI) error {
	if(ci.SCM == "gogs") {
		log.Info("gogs is worked")
	} else {
		log.Info("gogs is skipped")
	}
	
	return nil
}

func (c *GogsPlugin) Notify() error {
	return nil
}