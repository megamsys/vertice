package app

import (
	"github.com/megamsys/libgo/db"
	log "code.google.com/p/log4go"
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

type AssemblyResult struct {
   Id             string    `json:"id"` 
   Name           string   `json:"name"` 
   Components     []*Component    
   policies       string   `json:"policies"`
   inputs         string    `json:"inputs"`
   operations     string    `json:"operations"` 
   CreatedAt      string   `json:"created_at"` 
   }

type Component struct {
	 Id                            string    `json:"id"` 
    Name                           string    `json:“name”`
    ToscaType                      string    `json:“tosca_type”`
    Requirements                  *ComponentRequirements  
    Inputs                        *CompomentInputs  
    ExternalManagementResource     string
    Artifacts                     *Artifacts 
    RelatedComponents              string
    Operations                    *ComponentOperations	
   	CreatedAt      		           string   `json:"created_at"` 
   }

type ComponentRequirements struct {
	Host                    string  `json:"host"`
	Dummy                   string  `json:"dummy"`
}

type CompomentInputs struct {
	Domain                    string  `json:"domain"`
	Port                      string  `json:"port"`
	UserName                  string  `json:"username"`
	Password                  string  `json:"password"`
	Version                   string  `json:"version"`
	Source                    string  `json:"source"`
	DesignInputs             *DesignInputs
	ServiceInputs            *ServiceInputs
}

type DesignInputs struct {
	Id                          string    `json:"id"` 
    X                           string    `json:“x”`
    Y                           string    `json:“y”`
    Z                           string    `json:“z”`
    Wires                       []string    `json:“wires”`
}

type ServiceInputs struct {
	DBName                          string    `json:"dbname"` 
    DBPassword                      string    `json:“dbpassword”`
}

type Artifacts struct {
	ArtifactType                 string    `json:"artifact_type"` 
    Content                      string    `json:“content”`
}

type ComponentOperations struct {
	OperationType                 string    `json:"operation_type"` 
    TargetResource                string    `json:“target_resource”`
}


func (com *Component) Get(comId string) error {
    log.Info("Get message %v", comId)
    conn, err := db.Conn("components")
	if err != nil {	
		return err
	}	
	//appout := &Requests{}
	ferr := conn.FetchStruct(comId, com)
	if ferr != nil {	
		return ferr
	}	
	defer conn.Close()
	
	
	return nil

}

func (asm *Assembly) Get(asmId string) (*AssemblyResult, error) {
    log.Info("Get Assembly message %v", asmId)
    var j = -1
    asmresult := &AssemblyResult{}
	conn, err := db.Conn("assembly")
	if err != nil {	
		return asmresult, err
	}	
	//appout := &Requests{}
	ferr := conn.FetchStruct(asmId, asm)
	if ferr != nil {	
		return asmresult, ferr
	}	
	log.Info("------------assemblycom-------------------")
	log.Info(asm)
	var arraycomponent = make([]*Component, len(asm.Components))
	for i := range asm.Components {
		if len(asm.Components[i]) > 1 {
		  componentID := asm.Components[i]
		  component := Component{Id: componentID }
          err := component.Get(componentID)
		  if err != nil {
		       log.Info("Error: Riak didn't cooperate:\n%s.", err)
		       return asmresult, err
		  }
		  log.Info("------------component-------------------")
	      log.Info(component)
	      j++
	      log.Info(j)
		  arraycomponent[j] = &component
		  }
	    }
	result := &AssemblyResult{Id: asm.Id, Name: asm.Name,  Components: arraycomponent, CreatedAt: asm.CreatedAt}
	log.Info("------------result-------------------")
	log.Info(result)
	defer conn.Close()
	
	
	return result, nil
}

