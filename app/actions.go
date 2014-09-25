package app

import (
    "encoding/json"
	"bytes"
	"errors"
	"fmt"
	"github.com/tsuru/config"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/exec"
	"github.com/megamsys/gulp/scm"
	"log"
	"text/template"
	"path"
	"bufio"
	"strings"
	"regexp"
	"os"
	"bitbucket.org/kardianos/osext"
)

type DRBDMaster struct { DRBD DRBDM `json:"drbd"` }
type DRBDM struct {
		Remotehost  string `json:"remote_host"`
		Sourcedir   string `json:"source_dir"`
		Master      bool   `json:"master"`
		Archive     string `json:"archive"`
	}

type DRBDSlave struct { DRBD DRBDS `json:"drbd"` }
type DRBDS struct {
		Remotehost  string  `json:"remote_host"`
		Sourcedir   string	`json:"source_dir"`
		Archive     string  `json:"archive"`
	   }

const (
	keyremote_repo = "remote_repo="
	keylocal_repo  = "local_repo="
	keyproject     = "project="
	kibana         ="kibana"
	kibanaTemplatePath = "conf/kibana"
	kibanaDashPath = "/var/www/kibana/app/dashboards"
	nginx_restart = "/etc/init.d/service nginx restart"
	nginx_stop = "/etc/init.d/service nginx stop"
	nginx_start = "/etc/init.d/service nginx start"
	rootPath  = "/tmp"
	defaultEnvPath = "conf/env.sh"
	drbd_mnt = "/drbd_mnt"
)

var ErrAppAlreadyExists = errors.New("there is already an app with this name.")

func CommandExecutor(app *App) (action.Result, error) {
	var e exec.OsExecutor
	var b bytes.Buffer
	var commandWords []string
	if (app.AppReqs!=nil) {
	commandWords = strings.Fields(app.AppReqs.LCApply)
	} else {
	commandWords = strings.Fields(app.AppConf.LCApply)
	}
	fmt.Println(commandWords, len(commandWords))

	if len(commandWords) > 0 {
		if err := e.Execute(commandWords[0], commandWords[1:], nil, &b, &b); err != nil {
			return nil, err
		}
	}

	log.Printf("%s", b)
	return &app, nil
}

//In this function to convert bytearray value to string.
func CToGoString(c []byte) string {
    n := -1
    for i, b := range c {
        if b == 0 {
            break
        }
        n = i
    }
    return string(c[:n+1])
}

//create a new file in rootpath and write the data into that file using bufio package.
func FileCreator(name string, json []byte) (action.Result, error) {
        filePath := path.Join(rootPath, name + ".json")
		JsonFile, err := filesystem().Create(filePath)

		if err != nil {
		   return nil, err
		}

		w := bufio.NewWriter(JsonFile)
		res, err := w.WriteString(CToGoString(json[:]))
		w.Flush()
       log.Printf("%s", res)
	   return &JsonFile, nil
}

// insertApp is an action that inserts an app in the database in Forward and
// removes it in the Backward.
//
// The first argument in the context must be an App or a pointer to an App.
var startApp = action.Action{
	Name: "startapp",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app App
		switch ctx.Params[0].(type) {
		case App:
			app = ctx.Params[0].(App)
		case *App:
			app = *ctx.Params[0].(*App)
		default:
			return nil, errors.New("First parameter must be App or *App.")
		}

		return CommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		app := ctx.FWResult.(*App)
		log.Printf("[%s] Nothing to recover for %s", app.Name)
	},
	MinParams: 1,
}


var stopApp = action.Action{
	Name: "stopapp",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app App
		switch ctx.Params[0].(type) {
		case App:
			app = ctx.Params[0].(App)
		case *App:
			app = *ctx.Params[0].(*App)
		default:
			return nil, errors.New("First parameter must be App or *App.")
		}

		return CommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		app := ctx.FWResult.(*App)
		log.Printf("[%s] Nothing to recover for %s", app.Name)
	},
	MinParams: 1,
}

var buildApp = action.Action{
	Name: "buildapp",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app App
		switch ctx.Params[0].(type) {
		case App:
			app = ctx.Params[0].(App)
		case *App:
			app = *ctx.Params[0].(*App)
		default:
			return nil, errors.New("First parameter must be App or *App.")
		}

		project, err := scm.Project()
		if err != nil {
			log.Printf("Could not find the project name in gulp.conf file: %s", err)
			return nil, errors.New("Could not find the project name in gulp.conf file")
		}

		builder, err := scm.Builder()
		if err != nil {
			log.Printf("Could not find the builder in gulp.conf file: %s", err)
			return nil, errors.New("Could not find the builder in gulp.conf file")

		}

		local_repo, err := scm.GetPath()
		if err != nil {
			log.Printf("Could not find the local repo  in gulp.conf file: %s", err)
			return nil, errors.New("Could not find the local repo in gulp.conf file")
		}

		remote_repo, err := scm.GetRemotePath()
		if err != nil {
			log.Printf("Could not find the remote repo in gulp.conf file: %s", err)
			return nil, errors.New("Could not find the remote repo in gulp.conf file")
		}

		build_parms := fmt.Sprintf("%s/%s %s %s %s", builder, app.AppReqs.LCApply, keyproject+project, keylocal_repo+local_repo, keyremote_repo+remote_repo)

		app.AppReqs.LCApply = build_parms
		return CommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
	app := ctx.FWResult.(*App)
		log.Printf("[%s] Nothing to recover for %s", app.Name)
	},
	MinParams: 1,
}

var launchedApp = action.Action{
	Name: "launchedapp",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app App
		switch ctx.Params[0].(type) {
		case App:
			app = ctx.Params[0].(App)
		case *App:
			app = *ctx.Params[0].(*App)
		default:
			return nil, errors.New("First parameter must be App or *App.")
		}

		log.Printf("Launched, attaching post install to %s", app.Name)

		tmpl, err := template.New(kibana).ParseFiles(kibanaTemplatePath)

		if err != nil {
	      return nil, err
		}

		kibanaPath := path.Join(kibanaDashPath, app.Name + ".json")
		kibanaFile, err := filesystem().Create(kibanaPath)

		if err != nil {
		   return nil, err
		}

		w := bufio.NewWriter(kibanaFile)
		err = tmpl.Execute(w, app)
		w.Flush()

		tmpAppreq := &AppRequests{}
		tmpAppreq.LCApply = nginx_restart

		app.AppReqs = tmpAppreq

		if err != nil {
	      return nil, err
		}

		return CommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		app := ctx.FWResult.(*App)
		log.Printf("[%s] Nothing to recover for %s", app.Name)
	},
	MinParams: 1,
}

var addonApp = action.Action{
	Name: "addonapp",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app App
		switch ctx.Params[0].(type) {
		case App:
			app = ctx.Params[0].(App)
		case *App:
			app = *ctx.Params[0].(*App)
		default:
			return nil, errors.New("First parameter must be App or *App.")
		}

		log.Printf("Addon, attaching post install to %s", app.Name)

	    if app.AppConf.DRFromhost != "" {
	    instance_name, err := config.GetString("name")
	    if err != nil {
		     return nil, err
	      }
	     localRepo, _ := config.GetString("scm:local_repo")
	    if instance_name == app.AppConf.DRFromhost {
	          group := DRBDMaster {
	                   DRBDM {
		                      Remotehost: app.AppConf.DRToHosts,
		                      Sourcedir:  localRepo,
		                      Master: true,
		                      Archive: app.AppConf.DRLocations,
	                         },
	                      }
	         b, err := json.Marshal(group)
	         if err != nil {
		           fmt.Println("error:", err)
		           return nil, err
	          }
	         log.Printf("Found Addon-DR, creating json %s", b)
	         FileCreator(instance_name, b)
	         } else {
	              group := DRBDSlave{
	                        DRBDS{
		                           Remotehost: app.AppConf.DRFromhost,
		                           Sourcedir: localRepo,
		                           Archive: app.AppConf.DRLocations,
	                            },
	               }
	           b, err := json.Marshal(group)
	          if err != nil {
		        fmt.Println("error:", err)
		        return nil, err
	          }
	         FileCreator(instance_name, b)
	        }
	      tmpAppConf := &AppConfigurations{}
		  tmpAppConf.LCApply = "chef-client -o '"+ app.AppConf.DRRecipe +"' -j /tmp/"+ instance_name + ".json"
	      app.AppConf = tmpAppConf
	     }
		return CommandExecutor(&app)
		},
		Backward: func(ctx action.BWContext) {
		app := ctx.FWResult.(*App)
		log.Printf("[%s] Nothing to recover for %s", app.Name)
	},
	MinParams: 1,
 }

var modifyEnv = action.Action{
	Name: "modifyEnv",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app App
		switch ctx.Params[0].(type) {
		case App:
			app = ctx.Params[0].(App)
		case *App:
			app = *ctx.Params[0].(*App)
		default:
			return nil, errors.New("First parameter must be App or *App.")
		}
		 var lines [][]string
         tmparray := make([]string, 10)
         var i int
		folderPath, err := osext.ExecutableFolder()
        if err != nil {
            return nil, err
         }
         Path := path.Join(folderPath + defaultEnvPath)
		 file, err := os.Open(Path)
         if err != nil {
            return nil, err
          }
         scanner := bufio.NewScanner(file)
         i = 0
         for scanner.Scan() {
            line := scanner.Text()
            if line != "" {
                    re, err := regexp.Compile(`MEGAM_APP_SERVICE_HOME=(.*)`)
                    if err != nil {
                        return nil, err
                    }
                    lines = re.FindAllStringSubmatch(line, -1)
                    i = i+1
                    if len(lines) > 0 {
                         tmparray[i] = strings.Replace(line, lines[0][1], drbd_mnt, 1)
                    } else {
                         tmparray[i] = line
                    }
              }
		 }
        defer file.Close()
        file1, err := filesystem().Create(Path)

		if err != nil {
		   return nil, err
		}
        w := bufio.NewWriter(file1)
         for i := range tmparray {
            res, err := w.WriteString(tmparray[i]+"\n")
            if err != nil {
		       return nil, err
		     }
            log.Printf("Change to HA location successful --> %s", res)
        }
		w.Flush()

	   return CommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		app := ctx.FWResult.(*App)
		log.Printf("[%s] Nothing to recover for %s", app.Name)
	},
	MinParams: 1,
	}

var nginxStart = action.Action{
	Name: "nginxStart",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app App
		switch ctx.Params[0].(type) {
		case App:
			app = ctx.Params[0].(App)
		case *App:
			app = *ctx.Params[0].(*App)
		default:
			return nil, errors.New("First parameter must be App or *App.")
		}
		tmpAppConf := &AppConfigurations{}
		//start the nginx server
		tmpAppConf.LCApply = nginx_start
	    app.AppConf = tmpAppConf
	   return CommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		app := ctx.FWResult.(*App)
		log.Printf("[%s] Nothing to recover for %s", app.Name)
	},
	MinParams: 1,
	}

  var nginxStop = action.Action{
	Name: "nginxStop",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app App
		switch ctx.Params[0].(type) {
		case App:
			app = ctx.Params[0].(App)
		case *App:
			app = *ctx.Params[0].(*App)
		default:
			return nil, errors.New("First parameter must be App or *App.")
		}
		tmpAppConf := &AppConfigurations{}
		//stop the nginx server
		tmpAppConf.LCApply = nginx_stop
	    app.AppConf = tmpAppConf
	   return CommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
	app := ctx.FWResult.(*App)
		log.Printf("[%s] Nothing to recover for %s", app.Name)
	},
	MinParams: 1,
	}

/*
// exportEnvironmentsAction exports megam's default environment variables in a
// new app. It requires a pointer to an App instance as the first parameter,
// and the previous result to be a *s3Env (it should be used after
var exportEnvironmentsAction = action.Action{
	Name: "export-environments",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		app := ctx.Params[0].(*App)
		err := app.Get()
		if err != nil {
			return nil, err
		}
		t, err := auth.CreateApplicationToken(app.Name)
		if err != nil {
			return nil, err
		}
		host, _ := config.GetString("host")
		envVars := []bind.EnvVar{
			{Name: "MEGAM_APPNAME", Value: app.Name},
			{Name: "MEGAM_HOST", Value: host},
			{Name: "MEGAM_API_KEY", Value: t.Token},
		}
		env, ok := ctx.Previous.(*s3Env)
		if ok {
			variables := map[string]string{
				"ENDPOINT":           env.endpoint,
				"LOCATIONCONSTRAINT": strconv.FormatBool(env.locationConstraint),
				"ACCESS_KEY_ID":      env.AccessKey,
				"SECRET_KEY":         env.SecretKey,
				"BUCKET":             env.bucket,
			}
			for name, value := range variables {
				envVars = append(envVars, bind.EnvVar{
					Name:         fmt.Sprintf("MEGAM_S3_%s", name),
					Value:        value,
					InstanceName: s3InstanceName,
				})
			}
		}
		err = app.setEnvsToApp(envVars, false, true)
		if err != nil {
			return nil, err
		}
		return ctx.Previous, nil
	},
	Backward: func(ctx action.BWContext) {
		app := ctx.Params[0].(*App)
		auth.DeleteToken(app.Env["MEGAM_API_KEY"].Value)
		if app.Get() == nil {
			s3Env := app.InstanceEnv(s3InstanceName)
			vars := make([]string, len(s3Env)+3)
			i := 0
			for k := range s3Env {
				vars[i] = k
				i++
			}
			vars[i] = "MEGAM_HOST"
			vars[i+1] = "MEGAM_APPNAME"
			vars[i+2] = "MEGAM_APIKEY"
			app.UnsetEnvs(vars, false)
		}
	},
	MinParams: 1,
}

*/
