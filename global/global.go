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
package global

import (
	"github.com/megamsys/libgo/db"
	"crypto/rand"
    "math/big"
    "strings"
	log "code.google.com/p/log4go"
)

type Message struct {
	Id string `json:"id"`
}

type EventMessage struct {
	AssemblyId string `json:"assembly_id"`
	ComponentId string `json:"component_id"`
	Event string `json:"event"`
}

type PredefClouds struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Accounts_id string     `json:"accounts_id"`
	Jsonclaz    string     `json:"json_claz"`
	Spec        *PDCSpec   `json:"spec"`
	Access      *PDCAccess `json:"access"`
	Ideal       string     `json:"ideal"`
	CreatedAT   string     `json:"created_at"`
	Performance string     `json:"performance"`
}

type PDCSpec struct {
	TypeName string `json:"type_name"`
	Groups   string `json:"groups"`
	Image    string `json:"image"`
	Flavor   string `json:"flavor"`
	TenantID string `json:"tenant_id"`
}

type PDCAccess struct {
	Sshkey         string `json:"ssh_key"`
	IdentityFile   string `json:"identity_file"`
	Sshuser        string `json:"ssh_user"`
	VaultLocation  string `json:"vault_location"`
	SshpubLocation string `json:"sshpub_location"`
	Zone           string `json:"zone"`
	Region         string `json:"region"`
}

type Component struct {
	Id                         string `json:"id"`
	Name                       string `json:"name"`
	ToscaType                  string `json:"tosca_type"`
	Requirements               *ComponentRequirements
	Inputs                     *CompomentInputs
	ExternalManagementResource string
	Artifacts                  *Artifacts
	RelatedComponents          string
	Operations                 *ComponentOperations
	CreatedAt                  string `json:"created_at"`
}

/**
**fetch the component json from riak and parse the json to struct
**/
func (asm *Component) Get(asmId string) (*Component, error) {
    log.Info("Get Component message %v", asmId)
    conn, err := db.Conn("components")
	if err != nil {	
		return asm, err
	}	
	//appout := &Requests{}
	ferr := conn.FetchStruct(asmId, asm)
	if ferr != nil {	
		return asm, ferr
	}	
	defer conn.Close()
	
	return asm, nil

}


type ComponentRequirements struct {
	Host  string `json:"host"`
	Dummy string `json:"dummy"`
}

type CompomentInputs struct {
	Domain        string `json:"domain"`
	Port          string `json:"port"`
	UserName      string `json:"username"`
	Password      string `json:"password"`
	Version       string `json:"version"`
	Source        string `json:"source"`
	DesignInputs  *DesignInputs
	ServiceInputs *ServiceInputs
	CIID          string  `json:"ci_id"`
}

type DesignInputs struct {
	Id    string   `json:"id"`
	X     string   `json:“x”`
	Y     string   `json:“y”`
	Z     string   `json:“z”`
	Wires []string `json:“wires”`
}

type ServiceInputs struct {
	DBName     string `json:"dbname"`
	DBPassword string `json:"dbpassword"`
}

type Artifacts struct {
	ArtifactType string `json:"artifact_type"`
	Content      string `json:"content"`
}

type ComponentOperations struct {
	OperationType  string `json:"operation_type"`
	TargetResource string `json:"target_resource"`
}

type Request struct {
	Id	             string     `json:"id"`
	NodeId           string     `json:"node_id"`
	NodeName         string     `json:"node_name"` 
	ReqType          string     `json:"req_type"`
	CreatedAt        string     `json:"created_at"`
}

/**
**fetch the request json from riak and parse the json to struct
**/
func (req *Request) Get(reqId string) (*Request, error) {
    log.Info("Get Request message %v", reqId)
    conn, err := db.Conn("requests")
	if err != nil {	
		return req, err
	}	
	//appout := &Requests{}
	ferr := conn.FetchStruct(reqId, req)
	if ferr != nil {	
		return req, ferr
	}	
	defer conn.Close()
	
	return req, nil

}

type Assemblies struct {
   Id             string    `json:"id"` 
   AccountsId     string    `json:"accounts_id"`
   JsonClaz       string   `json:"json_claz"` 
   Name           string   `json:"name"` 
   Assemblies     []string   `json:"assemblies"` 
   Inputs         *AssembliesInputs   `json:"inputs"` 
   CreatedAt      string   `json:"created_at"` 
   }

type AssembliesInputs struct {
   Id                   string    `json:"id"` 
   AssembliesType       string    `json:"assemblies_type"` 
   Label                string    `json:"label"` 
   CloudSettings        []*CloudSettings    `json:"cloudsettings"`
   }

type CloudSettings struct {
	Id                 string       `json:"id"`
    CSType             string        `json:"cstype"`
    CloudSettings      string       `json:"cloudsettings"`
    X                  string        `json:"x"`
    Y                  string        `json:"y"`
    Z                  string        `json:"z"`
    Wires              []string    `json:“wires”`
}

type CI struct {
	Enable				string		`json:"enable"`
	SCM					string		`json:"scm"`
	Token				string		`json:"token"`
	Owner				string		`json:"owner"`
	ComponentID			string		`json:"component_id"`
	AssemblyID			string		`json:"assembly_id"`
	Id					string		`json:"id"`
	CreatedAT			string		`json:"created_at"`
}

/**
**fetch the continious integration data from riak and parse the json to struct
**/
func (req *CI) Get(reqId string) (*CI, error) {
    log.Info("Get Continious Integration message %v", reqId)
    conn, err := db.Conn("cig")
	if err != nil {	
		return req, err
	}	
	//appout := &Requests{}
	ferr := conn.FetchStruct(reqId, req)
	if ferr != nil {	
		return req, ferr
	}	
	defer conn.Close()
	
	return req, nil

}

type Output struct {
	Key     string   `json:"key"`
	Value   string   `json:"value"`
}

type Policy struct {
	Name    string   `json:"name"`
	Ptype   string   `json:"ptype"`
	Members []string `json:"members"`
}

type Assembly struct {
   Id             string   	 	`json:"id"` 
   JsonClaz       string   		`json:"json_claz"` 
   Name           string   		`json:"name"` 
   ToscaType      string        `json:"tosca_type"`
   Components     []string   	`json:"components"` 
   Policies       []*Policy  	`json:"policies"`
   Inputs         []*Output    	`json:"inputs"`
   Operations     string    	`json:"operations"` 
   Outputs        []*Output  	`json:"outputs"`
   Status         string    	`json:"status"`
   CreatedAt      string   		`json:"created_at"` 
   }
   
type AssemblyResult struct {
	Id         string 			`json:"id"`
	Name       string 			`json:"name"`
	ToscaType  string           `json:tosca_type"`
	Components []*Component		
	policies   []*Policy 		`json:"policies"`
	inputs     string    		`json:"inputs"`
	operations string    		`json:"operations"`
	Command    string
	CreatedAt  string 			`json:"created_at"`
}
   
/**
**fetch the Assembly data from riak and parse the json to struct
**/
func (req *Assembly) Get(reqId string) (*Assembly, error) {
    log.Info("Get Assembly message %v", reqId)
    conn, err := db.Conn("assembly")
	if err != nil {	
		return req, err
	}	
	//appout := &Requests{}
	ferr := conn.FetchStruct(reqId, req)
	if ferr != nil {	
		return req, ferr
	}	
	defer conn.Close()
	
	return req, nil

}

func (asm *Assembly) GetResult(asmId string) (*AssemblyResult, error) {
    log.Info("Get Assembly message %v", asmId)
    var j = -1
    asmresult := &AssemblyResult{}
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
	var arraycomponent = make([]*Component, len(asm.Components))
	for i := range asm.Components {
		 log.Debug(asm.Components[i])
		 t := strings.TrimSpace(asm.Components[i])
		 log.Debug(t)
		 log.Debug(len(t))
		if len(t) > 1  {
			log.Debug("entry")
		  componentID := asm.Components[i]
		  component := Component{Id: componentID }
          com, err := component.Get(componentID)
		  if err != nil {
		       log.Error("Error: Riak didn't cooperate:\n%s.", err)
		       return asmresult, err
		  }
	      j++
	      log.Debug(j)
	      log.Debug(asm.Components[i])
		  arraycomponent[j] = com
		  }
	    }
	log.Info("else entry")
	result := &AssemblyResult{Id: asm.Id, Name: asm.Name, ToscaType: asm.ToscaType,  Components: arraycomponent, CreatedAt: asm.CreatedAt}
	defer conn.Close()
	
	return result, nil
}

/**
generate the rand string 
**/
func RandString(n int) string {
    const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    symbols := big.NewInt(int64(len(alphanum)))
    states := big.NewInt(0)
    states.Exp(symbols, big.NewInt(int64(n)), nil)
    r, err := rand.Int(rand.Reader, states)
    if err != nil {
        panic(err)
    }
    var bytes = make([]byte, n)
    r2 := big.NewInt(0)
    symbol := big.NewInt(0)
    for i := range bytes {
        r2.DivMod(r, symbols, symbol)
        r, r2 = r2, r
        bytes[i] = alphanum[symbol.Int64()]
    }
    return string(bytes)
}

