package app

import (
	"encoding/json"
	"log"
	"github.com/megamsys/libgo/fs"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/db"
	"github.com/megamsys/gulp/app/bind"
	"regexp"
)

var (
	cnameRegexp = regexp.MustCompile(`^[a-zA-Z0-9][\w-.]+$`)
	fsystem fs.Fs
)

// Appreq is the main type in megam. An app represents a real world application.
// This struct holds information about the app: its name, address, list of
// teams that have access to it, used platform, etc.
type App struct {
	Env      map[string]bind.EnvVar
	Id	     string
	Platform string `chef:"java"`
	Name     string
	Ip       string
	Type     string
	CName    string
	//	Units    []Unit
	State   string
	Deploys uint
    AppReqs *AppRequests
    AppConf *AppConfigurations
	//	hr hookRunner
}

type AppRequests struct {
   AppId             string    `json:"id"`
   NodeId         string   `json:"node_id"`
   NodeName       string   `json:"node_name"`
   AppDefnsId     string   `json:"appdefns_id"`
   ReqType        string   `json:"req_type"`
   LCApply        string   `json:"lc_apply"`
   LCAdditional   string   `json:"lc_additional"`
   LCWhen         string   `json:"lc_when"`
   CreatedAT      string   `json:"created_at"`
   }

   type AppConfigurations struct {
   		ConfigId       		string   `json:"id"`
   		NodeId         		string   `json:"node_id"`
   		NodeName       		string   `json:"node_name"`
   		DRLocations    		string
   		DRFromhost     		string
   		DRToHosts      		string
   		DRRecipe            string
   		HAProxyhost    		string
   		LoadbalancedHosts 	string
   		LoadRecipe            string
   		CPUThreshhold    	string
   		MemThreshhold    	string
   		Noofinstances       string
   		AutoRecipe            string
   		MonitoringAgent     string
   		MonitorRecipe            string
   		CreatedAT      		string   `json:"created_at"`
   		LCApply             string
   }

// MarshalJSON marshals the app in json format. It returns a JSON object with
//the following keys: name, framework, teams, units, repository and ip.
func (app *App) MarshalJSON() ([]byte, error) {
	result := make(map[string]interface{})
	result["name"] = app.Name
	result["platform"] = app.Platform
	//result["repository"] = repository.ReadWriteURL(app.Name)
	result["ip"] = app.Ip
	result["cname"] = app.CName
	result["launched"] = app.State == "launched"
	return json.Marshal(&result)
}

//UnmarshalJSON parse app configuration json using AppConfiguartiosn struct.
func (a *AppConfigurations) UnmarshalJSON(b []byte) error {

	var f interface{}
	json.Unmarshal(b, &f)

	m := f.(map[string]interface{})
    a.ConfigId = m["id"].(string)
    a.NodeId  = m["node_id"].(string)
    a.NodeName  = m["node_name"].(string)

	config := m["config"]
	conf := config.(map[string]interface{})

    dis := conf["disaster"]
    disaster := dis.(map[string]interface{})
    a.DRLocations = disaster["locations"].(string)
    a.DRFromhost  = disaster["fromhost"].(string)
    a.DRToHosts  = disaster["tohosts"].(string)
    a.DRRecipe   = disaster["recipe"].(string)

    load := conf["loadbalancing"]
    loadbalance := load.(map[string]interface{})
    a.HAProxyhost = loadbalance["haproxyhost"].(string)
    a.LoadbalancedHosts  = loadbalance["loadbalancehost"].(string)
    a.LoadRecipe   = loadbalance["recipe"].(string)

    scale := conf["autoscaling"]
    autoscale := scale.(map[string]interface{})
    a.CPUThreshhold  = autoscale["cputhreshold"].(string)
    a.MemThreshhold  = autoscale["memorythreshold"].(string)
    a.Noofinstances  = autoscale["noofinstances"].(string)
    a.AutoRecipe     = autoscale["recipe"].(string)

    mon := conf["monitoring"]
    monitor := mon.(map[string]interface{})
    a.MonitoringAgent  = monitor["agent"].(string)
    a.MonitorRecipe   = monitor["recipe"].(string)
    return nil
}

func filesystem() fs.Fs {
	if fsystem == nil {
		fsystem = fs.OsFs{}
	}
	return fsystem
}

// Get queries the database and fills the App object with data retrieved from
// the database. It uses the name of the app as filter in the query, so you can
// provide this field:
//
//     app := App{Name: "myapp"}
//     err := app.Get()
//     // do something with the app
func (app *App) Get(reqId string) error {
log.Printf("Get message %v", reqId)
	if app.Type != "addon" {
	conn, err := db.Conn("appreqs")
	if err != nil {
		return err
	}
	appout := &AppRequests{}
	conn.FetchStruct(reqId, appout)
	app.AppReqs = appout
	defer conn.Close()
	} else {
	  conn, err := db.Conn("addonconfigs")
	if err != nil {
		return err
	}
	appout := &AppConfigurations{}
	conn.FetchStruct(reqId, appout)
	app.AppConf = appout
	log.Printf("Get message from riak  %v", appout)
	defer conn.Close()
	}
	//fetch it from riak.
	// conn.Fetch(app.id)
	// store stuff back in the appreq object.
	return nil
}

// StartsApp creates a new app.
//
// Starts the app :
func StartApp(app *App) error {
	actions := []*action.Action{&startApp}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}

// StopsApp creates a new app.
//
// Stops the app :
func StopApp(app *App) error {
	actions := []*action.Action{&stopApp}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}

// StopsApp creates a new app.
//
// Stops the app :
func BuildApp(app *App) error {
	actions := []*action.Action{&buildApp}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}

// StopsApp creates a new app.
//
// Stops the app :
func LaunchedApp(app *App) error {
	actions := []*action.Action{&launchedApp}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}

//Addon action for App
func AddonApp(app *App) error {
    actions := []*action.Action{&stopApp, &addonApp, &modifyEnv, &startApp}

    pipeline := action.NewPipeline(actions...)
    err := pipeline.Execute(app)
    if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}

// GetName returns the name of the app.
func (app *App) GetName() string {
	return app.Name
}

// GetIp returns the ip of the app.
func (app *App) GetIp() string {
	return app.Ip
}

// GetIp returns the ip of the app.
func (app *App) GetType() string {
	return app.Type
}

// GetPlatform returns the platform of the app.
func (app *App) GetPlatform() string {
	return app.Platform
}

func (app *App) GetDeploys() uint {
	return app.Deploys
}

// Env returns app.Env
func (app *App) Envs() map[string]bind.EnvVar {
	return app.Env
}

// GetAppReqs returns the app requests of the app.
func (app *App) GetAppReqs() *AppRequests {
	return app.AppReqs
}

func (app *App) GetAppConf() *AppConfigurations {
    return app.AppConf
}

/* setEnv sets the given environment variable in the app.
func (app *App) setEnv(env bind.EnvVar) {
	if app.Env == nil {
		app.Env = make(map[string]bind.EnvVar)
	}
	app.Env[env.Name] = env
	if env.Public {
		app.Log(fmt.Sprintf("setting env %s with value %s", env.Name, env.Value), "megam")
	}
}

// getEnv returns the environment variable if it's declared in the app. It will
// return an error if the variable is not defined in this app.
func (app *App) getEnv(name string) (bind.EnvVar, error) {
	var (
		env bind.EnvVar
		err error
		ok  bool
	)
	if env, ok = app.Env[name]; !ok {
		err = stderr.New("Environment variable not declared for this app.")
	}
	return env, err
}
*/
