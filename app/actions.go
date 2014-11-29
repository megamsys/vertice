package app

import (
	"errors"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/exec"
	log "code.google.com/p/log4go"
	"strings"
	"github.com/tsuru/config"
	"os"
	"path"
	"bufio"
	"github.com/megamsys/megamd/provisioner"
	"github.com/megamsys/libgo/db"
)


func CommandExecutor(app *provisioner.AssemblyResult) (action.Result, error) {
    var e exec.OsExecutor
    var commandWords []string

    commandWords = strings.Fields(app.Command)
    log.Debug("Command Executor entry: %s\n", app)
    megam_home, ckberr := config.GetString("MEGAM_HOME")
	if ckberr != nil {
		return nil, ckberr
	}
    appName := app.Name + "." + app.Components[0].Inputs.Domain
	basePath := megam_home + "logs" 
	dir := path.Join(basePath, appName)
	
	fileOutPath := path.Join(dir, appName + "_out" )
	fileErrPath := path.Join(dir, appName + "_err" )
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Info("Creating directory: %s\n", dir)
		if errm := os.MkdirAll(dir, 0777); errm != nil {
			return nil, errm
		}
	} 
		// open output file
		fout, outerr := os.OpenFile(fileOutPath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if outerr != nil {
			return nil, outerr
		}
		defer fout.Close()
		// open Error file
		ferr, errerr := os.OpenFile(fileErrPath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if errerr != nil {
			return nil, errerr
		}
		defer ferr.Close()
  
	foutwriter := bufio.NewWriterSize(fout, 1)
	ferrwriter := bufio.NewWriterSize(ferr, 1)
    log.Debug(commandWords)
    log.Debug("Length: %s", len(commandWords))
    
    defer ferrwriter.Flush()
    defer foutwriter.Flush()
    
    if len(commandWords) > 0 {
       if err := e.Execute(commandWords[0], commandWords[1:], nil, foutwriter, ferrwriter); err != nil {
           return nil, err
        }
     }

  
  return &app, nil
}


var launchedApp = action.Action{
	Name: "launchedapp",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app provisioner.AssemblyResult
		switch ctx.Params[0].(type) {
		case provisioner.AssemblyResult:
			app = ctx.Params[0].(provisioner.AssemblyResult)
		case *provisioner.AssemblyResult:
			app = *ctx.Params[0].(*provisioner.AssemblyResult)
		default:
			return nil, errors.New("First parameter must be App or *assemblies.AssemblyResult.")
		}
		return CommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}

var updateStatus = action.Action{
	Name: "updatestatus",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app provisioner.AssemblyResult
		switch ctx.Params[0].(type) {
		case provisioner.AssemblyResult:
			app = ctx.Params[0].(provisioner.AssemblyResult)
		case *provisioner.AssemblyResult:
			app = *ctx.Params[0].(*provisioner.AssemblyResult)
		default:
			return nil, errors.New("First parameter must be App or *assemblies.AssemblyResult.")
		}
		asm := &Assembly{}
	    conn, err := db.Conn("assembly")
	     if err != nil {	
		    return nil, err
	      }	

	    ferr := conn.FetchStruct(app.Id, asm)
	    if ferr != nil {	
		   return nil, ferr
	     }	
	
	   update := Assembly{
		Id:            asm.Id, 
        JsonClaz:      asm.JsonClaz, 
        Name:          asm.Name, 
        Components:    asm.Components ,
        Policies:      asm.Policies,
        Inputs:        asm.Inputs,
        Operations:    asm.Operations,
        Outputs:       asm.Outputs,
        Status:        "Terminated",
        CreatedAt:     asm.CreatedAt,
	   }
	   err = conn.StoreStruct(app.Id, &update)		
		
		return CommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		log.Info("[%s] Nothing to recover")
	},
	MinParams: 1,
}




