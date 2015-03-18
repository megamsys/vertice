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

/**
** watching the application for CMP   
** get trigger url from config file 
**/
func (c *CMPPlugin) Watcher(ci *global.CI) error {
	if(ci.SCM == "CMP") {
		log.Info("CMP is worked")
	} else {
		log.Info("CMP is skipped")
	}
	
	return nil
}

/**
**notify the messages or any other operations to CMP
**/
func (c *CMPPlugin) Notify(m *global.EventMessage) error {
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
	if(ci.SCM == "CMP") {
		log.Info("CMP is worked")
	} else {
		log.Info("CMP is skipped")
	}
	return nil
}