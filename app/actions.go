package app

import (
	"bytes"
	"fmt"
	"errors"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/exec"
	log "code.google.com/p/log4go"
	"strings"
	"github.com/megamsys/megamd/provisioner"
)

func CommandExecutor(app *provisioner.AssemblyResult) (action.Result, error) {
	var e exec.OsExecutor
	var b bytes.Buffer
	
	commandWords := strings.Fields(app.Command)
	
	fmt.Println(commandWords, len(commandWords))
    fmt.Println("--------------------")
    fmt.Println(commandWords[0])
    fmt.Println(commandWords[4])
    s := make([]string, 1, 25)
    str := strings.Join(append(s[:2],commandWords...)," ")
    fmt.Println("=========================")
    fmt.Println(str)
    s1 := make([]string, 2)
    s1[0] = "-la"
	if len(commandWords) > 0 {
		//if err := e.Execute(commandWords[0], commandWords[1:len(commandWords)], nil, &b, &b); err != nil {
		if err := e.Execute("ls", s1, nil, &b, &b); err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	log.Info("%s", b)
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


