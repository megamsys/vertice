package coordinator

import (

	log "code.google.com/p/log4go"
	"github.com/megamsys/megamd/app"
)

type Coordinator struct {
	//RequestHandler(f func(*Message), name ...string) (Handler, error)
	//EventsHandler(f func(*Message), name ...string) (Handler, error)
}

func NewCoordinator(chann []byte, queue string) {
	log.Info("Handling coordinator message %v", string(chann))
	
	switch queue {
	case "Requests":
	      requestHandler(chann)
	      break;
	case "Events":
	      eventsHandler(chann)
	      break;      
}
}
	
func requestHandler(chann []byte) {
	    request := app.Request{Id: string(chann)}
        req, err := request.Get(string(chann))
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return
		}
	   log.Info("------------return request-------------------")
		       log.Info(req)
	   switch req.ReqType {
	   case "create":
	       	   assemblies := app.Assemblies{Id: req.AssembliesId }
               asm, err := assemblies.Get(req.AssembliesId)
		       if err != nil {
			         log.Error("Error: Riak didn't cooperate:\n%s.", err)
			         return
		         }
		       log.Info("------------assemblies-------------------")
		       log.Info(asm)
		       for i := range asm.Assemblies {
		       	if len(asm.Assemblies[i]) > 1 {
		             assemblyID := asm.Assemblies[i]
		             assembly := app.Assembly{Id: assemblyID }
                     res, err := assembly.Get(assemblyID)
		             if err != nil {
			            log.Error("Error: Riak didn't cooperate:\n%s.", err)
			            return
		              }
		             log.Info("------------assembly-------------------")
		             log.Info(res)
		             
		             //go CreateApp(&assembly)
	             }
		       	}	
		}
}	
	   
func eventsHandler(chann []byte) {
	
	
}	   
	