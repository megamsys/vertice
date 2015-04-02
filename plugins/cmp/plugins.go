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