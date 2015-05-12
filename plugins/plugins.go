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
package plugins

import (
  "github.com/megamsys/megamd/global"
	"fmt"
)


// Every Plugins must implement this interface.
type Plugins interface {
	// Called when watching a Machine.
	Watcher(*global.AssemblyWithComponents, *global.Operations, *global.Component) error

	// Called when notifing a Machine.
	Notify(*global.EventMessage) error
}


var plugs = make(map[string]Plugins)
var plug_names = []string{"cmp", "github", "gogs"}

/**
**register the all plugins to "plug" array
**/
func RegisterPlugins(name string, plugin Plugins) {
	plugs[name] = plugin
}

func GetPlugin(name string) (Plugins, error) {
	plugin, ok := plugs[name]
	if !ok {
		return nil, fmt.Errorf("plugins not registered")
	}
	return plugin, nil
}


func Watcher(asm *global.AssemblyWithComponents) error {
  if len(asm.Components) > 0 {
     	for i := range asm.Components {
     	  if asm.Components[i] != nil {
     		for j := range asm.Components[i].Operations {     			
				for k := range plug_names {
  					p, err := GetPlugin(plug_names[k])
	   				if err != nil {	
	      					return err
	  				 }		
					go p.Watcher(asm, asm.Components[i].Operations[j], asm.Components[i])  				
     		    } 
     		  }      		
     	  }
     }	  
  }   
  return nil
}

func Notify(m *global.EventMessage) error {
	for i := range plug_names {
  	p, err := GetPlugin(plug_names[i])
	   if err != nil {	
	      return err
	   }	
	perr :=  p.Notify(m)
	   if perr != nil {	
	      return perr
	   }	
  }
  return nil	
}


