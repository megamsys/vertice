package app

import (
	"github.com/megamsys/libgo/db"
	log "code.google.com/p/log4go"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/megamd/provisioner"
	"strings"
	"encoding/json"
	"github.com/tsuru/config"
	"github.com/megamsys/megamd/global"
)



type Assembly struct {
   Id             string    `json:"id"` 
   JsonClaz       string   `json:"json_claz"` 
   Name           string   `json:"name"` 
   Components     []string   `json:"components"` 
   Policies       []*provisioner.Policy   `json:"policies"`
   inputs         string    `json:"inputs"`
   operations     string    `json:"operations"` 
   CreatedAt      string   `json:"created_at"` 
   }



func (asm *Assembly) Get(asmId string) (*provisioner.AssemblyResult, error) {
    log.Info("Get Assembly message %v", asmId)
    var j = -1
    asmresult := &provisioner.AssemblyResult{}
    log.Debug(asmresult)
	log.Debug("--------asmresult-------")
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
		 log.Debug(asm.Components[i])
		 t := strings.TrimSpace(asm.Components[i])
		 log.Debug(t)
		 log.Debug(len(t))
		if len(t) > 1  {
			log.Debug("entry")
		  componentID := asm.Components[i]
		  component := provisioner.Component{Id: componentID }
          err := component.Get(componentID)
		  if err != nil {
		       log.Error("Error: Riak didn't cooperate:\n%s.", err)
		       return asmresult, err
		  }
	      j++
	      log.Debug(j)
	      log.Debug(asm.Components[i])
		  arraycomponent[j] = &component
		  }
	    }
	log.Info("else entry")
	result := &provisioner.AssemblyResult{Id: asm.Id, Name: asm.Name,  Components: arraycomponent, CreatedAt: asm.CreatedAt}
	defer conn.Close()
	
	
	return result, nil
}

func GetPredefClouds(host string) (*global.PredefClouds, error) {
	pdc := &global.PredefClouds{}

	predefBucket, perr := config.GetString("buckets:PREDEFCLOUDS")
	if perr != nil {
		return pdc, perr
	}
	conn, err := db.Conn(predefBucket)

	if err != nil {
		return pdc, err
	}

	ferr := conn.FetchStruct(host, pdc)
	if ferr != nil {
		return pdc, ferr
	}
	return pdc, nil
}

func LaunchApp(asm *provisioner.AssemblyResult, id string) error {
    log.Debug("Launch App entry")
	    com := &provisioner.Component{}
	    mapB, _ := json.Marshal(asm.Components[0])
        json.Unmarshal([]byte(string(mapB)), com)
        if com.Name != "" {
            s1, _ := GetPredefClouds(com.Requirements.Host)
           //	s := strings.Split(com.ToscaType, ".")
        	if s1.Spec.TypeName == "docker" {
        		log.Debug("Docker provisiner entry")
        		// Provisioner
	            p, err := provisioner.GetProvisioner("docker")
	            if err != nil {	
	                return err
	             }	
	            log.Info("Provisioner: %v", p)
	             _, perr := p.CreateCommand(asm, id)
	            if perr != nil {	
	               return perr
	             }	       		
        	} else {
        		// Provisioner
	            p, err := provisioner.GetProvisioner("chef")
	            if err != nil {	
	                return err
	             }	
	            log.Info(p)
	
	            str, perr := p.CreateCommand(asm, id)
	            if perr != nil {	
	               return perr
	             }	
	            asm.Command = str
	            actions := []*action.Action{&launchedApp}

	            pipeline := action.NewPipeline(actions...)
	            aerr := pipeline.Execute(asm)
	            if aerr != nil {
		           return aerr
	             } 
        	}
          }
    return nil
}

func DeleteApp(asm *provisioner.AssemblyResult, id string) error {
    log.Debug("Delete App entry")
	    com := &provisioner.Component{}
	    mapB, _ := json.Marshal(asm.Components[0])
        json.Unmarshal([]byte(string(mapB)), com)
        if com.Name != "" {
            s1, _ := GetPredefClouds(com.Requirements.Host)
           //	s := strings.Split(com.ToscaType, ".")
        	if s1.Spec.TypeName == "docker" {
        		log.Debug("Docker provisioner entry")
        		// Provisioner
	            p, err := provisioner.GetProvisioner("docker")
	            if err != nil {	
	                return err
	             }	
	            log.Info("Provisioner: %v", p)
	             _, perr := p.CreateCommand(asm, id)
	            if perr != nil {	
	               return perr
	             }	       		
        	} else {
        		// Provisioner
	            p, err := provisioner.GetProvisioner("chef")
	            if err != nil {	
	                return err
	             }	
	            log.Info(p)
	
	            str, perr := p.DeleteCommand(asm, id)
	            if perr != nil {	
	               return perr
	             }	
	            asm.Command = str
	            actions := []*action.Action{&launchedApp}

	            pipeline := action.NewPipeline(actions...)
	            aerr := pipeline.Execute(asm)
	            log.Debug(aerr)
	            log.Debug("---------")
	            if aerr != nil {
		           return aerr
	             } 
        	}
          }
    return nil
}


