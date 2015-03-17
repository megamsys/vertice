package plugins

import (
  //log "code.google.com/p/log4go"
  "github.com/megamsys/megamd/global"
	"fmt"
)


// Every Plugins must implement this interface.
type Plugins interface {
	// Called when watching a Machine.
	Watcher(*global.CI) error

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
	//return nil
}


func Watcher(ci *global.CI) error {
  for i := range plug_names {
  	p, err := GetPlugin(plug_names[i])
	   if err != nil {	
	      return err
	   }	
	perr :=  p.Watcher(ci)
	   if perr != nil {	
	      return perr
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


