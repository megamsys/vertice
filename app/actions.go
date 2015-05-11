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
	"github.com/megamsys/libgo/db"
	"github.com/megamsys/megamd/global"
)


func CommandExecutor(app *global.AssemblyWithComponents) (action.Result, error) {
    var e exec.OsExecutor
    var commandWords []string
    appName := ""
    commandWords = strings.Fields(app.Command)
    log.Debug("Command Executor entry: %s\n", app)
    megam_home, ckberr := config.GetString("MEGAM_HOME")
	if ckberr != nil {
		return nil, ckberr
	}
	pair, perr := global.ParseKeyValuePair(app.Inputs, "domain")
		if perr != nil {
			log.Error("Failed to get the domain value : %s", perr)
		}
	if len(app.Components) > 0 {
		appName = app.Name + "." + pair.Value
	} else {
		domainkey, err_domainkey := config.GetString("DOMAIN")
		if err_domainkey != nil {
			return nil, err_domainkey
		}
		appName = app.Name + "." + domainkey
	}    
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
		var app global.AssemblyWithComponents
		switch ctx.Params[0].(type) {
		case global.AssemblyWithComponents:
			app = ctx.Params[0].(global.AssemblyWithComponents)
		case *global.AssemblyWithComponents:
			app = *ctx.Params[0].(*global.AssemblyWithComponents)
		default:
			return nil, errors.New("First parameter must be App or *global.AssemblyWithComponents.")
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
		var app global.AssemblyWithComponents
		switch ctx.Params[0].(type) {
		case global.AssemblyWithComponents:
			app = ctx.Params[0].(global.AssemblyWithComponents)
		case *global.AssemblyWithComponents:
			app = *ctx.Params[0].(*global.AssemblyWithComponents)
		default:
			return nil, errors.New("First parameter must be App or *global.AssemblyWithComponents.")
		}
		asm := &global.Assembly{}
	    conn, err := db.Conn("assembly")
	     if err != nil {	
		    return nil, err
	      }	

	    ferr := conn.FetchStruct(app.Id, asm)
	    if ferr != nil {	
		   return nil, ferr
	     }	
	
	   update := global.Assembly{
		Id:            asm.Id, 
        JsonClaz:      asm.JsonClaz, 
        Name:          asm.Name, 
        Components:    asm.Components,
        Requirements:  asm.Requirements,
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




