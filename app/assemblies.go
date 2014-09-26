package app

import (
	"github.com/megamsys/megamd/app/bind"
	log "code.google.com/p/log4go"
	"github.com/megamsys/libgo/db"
)


type Request struct {
	Env              map[string]bind.EnvVar
	Id	             string     `json:"id"`
	AssembliesId     string   `json:"assemblies_id"` 
	ReqType           string     `json:"req_type"`
	CreatedAt        string     `json:"created_at"`
}

type Assemblies struct {
   Id             string    `json:"id"` 
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
   }


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
	
	
	return asm, nil

}