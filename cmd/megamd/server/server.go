package server

import (
	"github.com/megamsys/megamd/api/http"
	"github.com/megamsys/libgo/db"
	"github.com/megamsys/megamd/cmd/megamd/server/queue"
	log "code.google.com/p/log4go"
	"github.com/megamsys/libgo/etcd"
	"github.com/megamsys/megamd/app"
	"github.com/tsuru/config"
	"encoding/json"
	"github.com/megamsys/libgo/amqp"
	"strings"
)

type Server struct {
	HttpApi      *http.HttpServer
	QueueServers []*queue.QueueServer
	//AdminServer  *admin.HttpServer
	stopped      bool
}

type Status struct {
	Id    string   `json:"id"`
	Status  string  `json:"status"`
}


func NewServer() (*Server, error) {

    log.Info("Starting New server")
	httpApi := http.NewHttpServer()

	return &Server{
		HttpApi:     httpApi,
		}, nil
}

func (self *Server) ListenAndServe() error {
	log.Info("Starting admin interface on port")
	//var etcdServerList [2]string
    var queueInput [2]string
	queueInput[0] = "cloudstandup"
	queueInput[1] = "Events"
	
	// Queue input
	for i := range queueInput {
		listenQueue := queueInput[i]
		queueserver := queue.NewServer(listenQueue)
		go queueserver.ListenAndServe()
	}
	self.Watcher()
	self.HttpApi.ListenAndServe()
	return nil
}

func (self *Server) Watcher() {
	path, _ := config.GetString("etcd:path")
	log.Info(path)
	c := etcd.NewClient(path+"/")
	
	ch := make(chan *etcd.Response, 1)
	stop := make(chan bool, 1)
	
	dir, _ := config.GetString("etcd:directory")
	c.CreateDir(dir)
   
	go receiver(ch, stop)
   
    _, err := c.Watch("/"+dir, 0, true, ch, stop)
   	if err != nil {
		log.Error(err)
	}
	if err != etcd.ErrWatchStoppedByUser {
		log.Error("Watch returned a non-user stop error")
	}

	
}

func (self *Server) Stop() {
	if self.stopped {
		return
	}
	log.Info("Stopping server")
	self.stopped = true

	log.Info("Stopping api server")
	//self.HttpApi.Close()
	log.Info("Api server stopped")

}

func receiver(c chan *etcd.Response, stop chan bool) {
     for {
        select {
        case msg := <-c:
           log.Info("==================>receiver entry")
		   log.Info(msg.Node.Key)  
		   handler(msg)
        }
     }
	
	stop <- true
}


func handler(msg *etcd.Response) {
	asm := &app.Assembly{}
	
    res := &Status{}
    json.Unmarshal([]byte(msg.Node.Value), &res)
    
	conn, err := db.Conn("assembly")
	if err != nil {	
		log.Error(err)
	}	
	
	ferr := conn.FetchStruct(res.Id, asm)
	if ferr != nil {	
		log.Error(ferr)
	}	
	
	for i := range asm.Policies {
	  mapD := map[string]string{"id": res.Id, "policy_name": asm.Policies[i].Name}
      mapB, _ := json.Marshal(mapD)
      log.Info(string(mapB))
      publisher(msg.Node.Key, string(mapB))
    }
}

func publisher(key string, json string) {
	factor, aerr := amqp.Factory()
	if aerr != nil {
		log.Error("Failed to get the queue instance: %s", aerr)
	}
	s := strings.Split(key, "/")
	pubsub, perr := factor.Get(s[len(s)-1])
	if perr != nil {
		log.Error("Failed to get the queue instance: %s", perr)
	}
	
	 serr := pubsub.Pub([]byte(json))
	 if serr != nil {
		log.Error("Failed to publish the queue instance: %s", serr)
	}
}

