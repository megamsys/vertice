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
    "errors"
	log "code.google.com/p/log4go"
)

type Message struct {
	Id string `json:"id"`
}

type EventMessage struct {
	AssemblyId     string 	`json:"assembly_id"`
	ComponentId    string 	`json:"component_id"`
	Event          string 	`json:"event"`
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

type Request struct {
	Id	             string     `json:"id"`
	AssembliesId     string     `json:"cat_id"`
	AssembliesName   string     `json:"name"` 
	ReqType          string     `json:"cattype"`
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

type KeyValuePair struct {
	Key     string   `json:"key"`
	Value   string   `json:"value"`
}

type Policy struct {
	Name    string   `json:"name"`
	Ptype   string   `json:"ptype"`
	Members []string `json:"members"`
}

type Operations struct {
	OperationType 				string 				`json:"operation_type"`
	Description 				string				`json:"description"`
	OperationRequirements		[]*KeyValuePair		`json:"operation_requirements"`
}

type Artifacts struct {
	ArtifactType 			string 			`json:"artifact_type"`
	Content     		 	string 			`json:"content"`
	ArtifactRequirements  	[]*KeyValuePair	`json:"artifact_requirements"`
}

type Component struct {
	Id                         string 				`json:"id"`
	Name                       string 				`json:"name"`
	ToscaType                  string 				`json:"tosca_type"`
	Inputs                     []*KeyValuePair		`json:"inputs"`
	Outputs					   []*KeyValuePair		`json:"outputs"`
	Artifacts                  *Artifacts			`json:"artifacts"`
	RelatedComponents          []string				`json:"related_components"`
	Operations     			   []*Operations    	`json:"operations"` 
	Status         			   string    			`json:"status"`
	CreatedAt                  string 				`json:"created_at"`
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

type Assemblies struct {
   Id             string  	    	`json:"id"` 
   AccountsId     string    		`json:"accounts_id"`
   JsonClaz       string   			`json:"json_claz"` 
   Name           string   			`json:"name"` 
   Assemblies     []string   		`json:"assemblies"` 
   Inputs         []*KeyValuePair   `json:"inputs"` 
   CreatedAt      string   			`json:"created_at"` 
   }
   
 func (asm *Assemblies) Get(asmId string) (*Assemblies, error) {
    log.Info("Get Assemblies message %v", asmId)
    conn, err := db.Conn("assemblies")
	if err != nil {	
		return asm, err
	}	
	//appout := &Requests{}
	ferr := conn.FetchStruct(asmId, asm)
	if ferr != nil {	
		return asm, ferr
	}	
	defer conn.Close()
	log.Debug(asm)
	log.Debug("----------ASSEMBLIES--------")
	return asm, nil

}


type Assembly struct {
   Id             string   	 		`json:"id"` 
   JsonClaz       string   			`json:"json_claz"` 
   Name           string   			`json:"name"` 
   ToscaType      string        	`json:"tosca_type"`
   Components     []string   		`json:"components"` 
   Requirements	  []*KeyValuePair	`json:"requirements"`
   Policies       []*Policy  		`json:"policies"`
   Inputs         []*KeyValuePair   `json:"inputs"`
   Operations     []*Operations    	`json:"operations"` 
   Outputs        []*KeyValuePair  	`json:"outputs"`
   Status         string    		`json:"status"`
   CreatedAt      string   			`json:"created_at"` 
   }
   
type AssemblyWithComponents struct {
	Id         		string 				`json:"id"`
	Name       		string 				`json:"name"`
	ToscaType  		string          	`json:tosca_type"`
	Components 		[]*Component		
	Requirements	[]*KeyValuePair		`json:"requirements"`
    Policies        []*Policy  			`json:"policies"`
    Inputs          []*KeyValuePair   	`json:"inputs"`
    Operations      []*Operations    	`json:"operations"` 
    Outputs         []*KeyValuePair  	`json:"outputs"`
    Status          string    			`json:"status"`
    Command         string
    CreatedAt       string   			`json:"created_at"` 
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

func (asm *Assembly) GetAssemblyWithComponents(asmId string) (*AssemblyWithComponents, error) {
    log.Info("Get Assembly message %v", asmId)
    var j = -1
    asmresult := &AssemblyWithComponents{}   
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
		 t := strings.TrimSpace(asm.Components[i])		
		if len(t) > 1  {
		  componentID := asm.Components[i]
		  component := Component{Id: componentID }
          com, err := component.Get(componentID)
		  if err != nil {
		       log.Error("Error: Riak didn't cooperate:\n%s.", err)
		       return asmresult, err
		  }
	      j++	     
		  arraycomponent[j] = com
		  }
	    }
	log.Info("else entry")
	result := &AssemblyWithComponents{Id: asm.Id, Name: asm.Name, ToscaType: asm.ToscaType,  Components: arraycomponent, Requirements: asm.Requirements, Policies: asm.Policies, Inputs: asm.Inputs, Outputs: asm.Outputs, Operations: asm.Operations, Status: asm.Status, CreatedAt: asm.CreatedAt}
	defer conn.Close()	
	return result, nil
}

func ParseKeyValuePair(keyvaluepair []*KeyValuePair, searchkey string) (*KeyValuePair, error) {
 	for i := range keyvaluepair {
		if keyvaluepair[i].Key == searchkey {
			return keyvaluepair[i], nil
		}
	}
	return nil, errors.New("The specific search key was not found in pair input...")
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

