package app

import (
	"github.com/megamsys/libgo/db"
	log "code.google.com/p/log4go"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/megamd/provisioner"
)



type Assembly struct {
   Id             string    `json:"id"` 
   JsonClaz       string   `json:"json_claz"` 
   Name           string   `json:"name"` 
   Components     []string   `json:"components"` 
   policies       string   `json:"policies"`
   inputs         string    `json:"inputs"`
   operations     string    `json:"operations"` 
   CreatedAt      string   `json:"created_at"` 
   }


func (asm *Assembly) Get(asmId string) (*provisioner.AssemblyResult, error) {
    log.Info("Get Assembly message %v", asmId)
    var j = -1
    asmresult := &provisioner.AssemblyResult{}
	conn, err := db.Conn("assembly")
	if err != nil {	
		return asmresult, err
	}	
	//appout := &Requests{}
	ferr := conn.FetchStruct(asmId, asm)
	if ferr != nil {	
		return asmresult, ferr
	}	
	var arraycomponent = make([]*provisioner.Component, len(asm.Components))
	for i := range asm.Components {
		if len(asm.Components[i]) > 1 {
		  componentID := asm.Components[i]
		  component := provisioner.Component{Id: componentID }
          err := component.Get(componentID)
		  if err != nil {
		       log.Info("Error: Riak didn't cooperate:\n%s.", err)
		       return asmresult, err
		  }
	      j++
	      log.Info(j)
		  arraycomponent[j] = &component
		  }
	    }
	result := &provisioner.AssemblyResult{Id: asm.Id, Name: asm.Name,  Components: arraycomponent, CreatedAt: asm.CreatedAt}
	defer conn.Close()
	
	
	return result, nil
}


func LaunchApp(assembly *provisioner.AssemblyResult) error {
    
	// Provisioner
	p, err := provisioner.GetProvisioner("chef")
	if err != nil {	
	    return err
	}	
	log.Info(p)
	
	str, perr := p.CreateCommand(assembly)
	if perr != nil {	
	    return perr
	}	
	assembly.Command = str
	actions := []*action.Action{&launchedApp}

	pipeline := action.NewPipeline(actions...)
	aerr := pipeline.Execute(assembly)
	if aerr != nil {
		return aerr
	}
	return nil
}


