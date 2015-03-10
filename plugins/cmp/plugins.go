package cmp

import (
  "github.com/megamsys/megamd/plugins"
  "github.com/megamsys/megamd/global"
  log "code.google.com/p/log4go"
)


func Init() {
	plugins.RegisterPlugins("cmp", &CMPPlugin{})
}

type CMPPlugin struct{}



func (c *CMPPlugin) Watcher(ci *global.CI) error {
	if(ci.SCM == "CMP") {
		log.Info("CMP is worked")
	} else {
		log.Info("CMP is skipped")
	}
	
	return nil
}

func (c *CMPPlugin) Notify() error {
	return nil
}