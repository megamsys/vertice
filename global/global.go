package global

import (
	"github.com/megamsys/libgo/db"
	"crypto/rand"
    "math/big"
	log "code.google.com/p/log4go"
)

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
	DBPassword string `json:“dbpassword”`
}

type Artifacts struct {
	ArtifactType string `json:"artifact_type"`
	Content      string `json:“content”`
}

type ComponentOperations struct {
	OperationType  string `json:"operation_type"`
	TargetResource string `json:“target_resource”`
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

