package server

import (
	log "code.google.com/p/log4go"
	"encoding/json"
	"github.com/megamsys/libgo/amqp"
	"github.com/megamsys/libgo/db"
	"github.com/megamsys/libgo/etcd"
	"github.com/megamsys/megamd/api/http"
	"github.com/megamsys/megamd/app"
	"github.com/megamsys/megamd/global"
	"github.com/megamsys/megamd/cmd/megamd/server/queue"
	"github.com/tsuru/config"
//	"strings"
	"time"
	"os"
	"fmt"
)

type Server struct {
	HttpApi      *http.HttpServer
	QueueServers []*queue.QueueServer
	stopped      bool
}

type Status struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

func NewServer() (*Server, error) {

	log.Info("Starting New server")
	httpApi := http.NewHttpServer()

	return &Server{
		HttpApi: httpApi,
	}, nil
}

func (self *Server) ListenAndServe() error {
	log.Info("Starting admin interface on port")
	//var etcdServerList [2]string
	var queueInput [2]string
	queueInput[0] = "cloudstandup"
	queueInput[1] = "Events"
    self.Checker()
	// Queue input
	for i := range queueInput {
		listenQueue := queueInput[i]
		queueserver := queue.NewServer(listenQueue)
		go queueserver.ListenAndServe()
	}
	self.EtcdWatcher()
	self.HttpApi.ListenAndServe()
	return nil
}

type Connection struct {
	Dial     string `json:"dial"`
}

func (self *Server) Checker() {
	log.Info("Dialing Rabbitmq.......")
	factor, err := amqp.Factory()
	if err != nil {
		log.Error("Failed to get the queue instance: %s", err)
	}
	
	conn, connerr := factor.Dial()
    log.Debug("connection %v", conn)
    log.Debug("connection error %v", connerr)
    if connerr != nil {
    	 fmt.Fprintf(os.Stderr, "Error: %v\n Please start Rabbitmq service.\n", connerr)
         os.Exit(1)
    }
    log.Info("Rabbitmq connected")
    
    log.Info("Dialing Riak.......")
 
	 rconn, rerr := db.Conn("connection")
	 if rerr != nil {
		 fmt.Fprintf(os.Stderr, "Error: %v\n Please start Riak service.\n", connerr)
         os.Exit(1)
	 }

	 data := "sampledata"
	 ferr := rconn.StoreObject("sampleobject", data)
	 if ferr != nil {
	 	 fmt.Fprintf(os.Stderr, "Error: %v\n Please start Riak service.\n", ferr)
         os.Exit(1)
	 }
	 defer rconn.Close()
    log.Info("Riak connected")
	
}

// Next returns a channel which will emit an Event as soon as one of interest occurs
func (self *Server) EtcdWatcher() {
	rootPrefix := "/"
	etcdPath, _ := config.GetString("etcd:path")

	c := etcd.NewClient(etcdPath + rootPrefix)
    conn, connerr := c.Dial("tcp", "127.0.0.1:4001")
    log.Debug("client %v", c)
    log.Debug("connection %v", conn)
    log.Debug("connection error %v", connerr)
    
    if conn != nil {
	   log.Info(" [x] Etcd client %s", etcdPath, rootPrefix)

	   dir, _ := config.GetString("etcd:directory")
	   log.Info(" [x] Etcd Directory %s", dir)

	   etreschan := make(chan *etcd.Response, 1)
	   stop := make(chan bool, 1)

	   go func() {
		for {
			select {
			case <-stop:
				return
			default:
			   _, err1 := c.CreateDir(dir)

	           if err1 != nil {
		           log.Error(err1)
	           }
				log.Info(" [x] Watching %s", rootPrefix+dir)
				_, err := c.Watch(rootPrefix+dir, 0, true, etreschan, stop)
				log.Info(" [x] Watched (%s)", rootPrefix+dir)

				if err != nil {
					log.Info(" [x] Watched Error (%s)", rootPrefix+dir)
					log.Error(err)
				//	return
				}
				if err != etcd.ErrWatchStoppedByUser {
					log.Error("Watch returned a non-user stop error")
					//return
				}
				log.Info(" [x] Sleep-Watch (%s)", rootPrefix+dir)

				time.Sleep(time.Second)

				log.Info(" [x] Slept-Watch (%s)", rootPrefix+dir)
				//self.EtcdWatcher()
			}
		}

	}()

	go receiverEtcd(etreschan, stop)
  } else {
  	 fmt.Fprintf(os.Stderr, "Error: %v\n Please start etcd deamon.\n", connerr)
         os.Exit(1)
  }
}

func receiverEtcd(c chan *etcd.Response, stop chan bool) {
	for {
		select {
		case msg := <-c:
			log.Info(" [x] Handing etcd response (%v)", msg)
			if msg != nil {
				handlerEtcd(msg)
			} else {
				log.Info(" [x] Nil - Handing etcd response (%v)", msg)
				return
			}
		}
	}
	stop <- false
}

func handlerEtcd(msg *etcd.Response) {
	log.Info(" [x] Really Handle etcd response (%s)", msg.Node.Key)

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
	
	comp := &global.Component{}
	
	conn1, err1 := db.Conn("components")
	if err1 != nil {
		log.Error(err1)
	}

	ferr1 := conn1.FetchStruct(asm.Components[0], comp)
	if ferr1 != nil {
		log.Error(ferr1)
	}
    
	for i := range asm.Policies {
		mapD := map[string]string{"Id": res.Id, "Action": asm.Policies[i].Name}
		mapB, _ := json.Marshal(mapD)
		log.Info(string(mapB))
		asmname := asm.Name+"."+comp.Inputs.Domain
		//asmname := asm.Name
		publisher(asmname, string(mapB))
	}
}

func publisher(key string, json string) {
	factor, aerr := amqp.Factory()
	if aerr != nil {
		log.Error("Failed to get the queue instance: %s", aerr)
	}
	//s := strings.Split(key, "/")
	//pubsub, perr := factor.Get(s[len(s)-1])
	pubsub, perr := factor.Get(key)
	if perr != nil {
		log.Error("Failed to get the queue instance: %s", perr)
	}

	serr := pubsub.Pub([]byte(json))
	if serr != nil {
		log.Error("Failed to publish the queue instance: %s", serr)
	}
}

func (self *Server) Stop() {
	if self.stopped {
		return
	}
	log.Info("Stopping servers ....")
	self.stopped = true

	log.Info("Stopping API server")
	//self.HttpApi.Close()
	log.Info("Stopped API server")

}
