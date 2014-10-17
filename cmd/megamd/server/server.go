package server

import (
	log "code.google.com/p/log4go"
	"encoding/json"
	"github.com/megamsys/libgo/amqp"
	"github.com/megamsys/libgo/db"
	"github.com/megamsys/libgo/etcd"
	"github.com/megamsys/megamd/api/http"
	"github.com/megamsys/megamd/app"
	"github.com/megamsys/megamd/cmd/megamd/server/queue"
	"github.com/tsuru/config"
	"strings"
	"time"
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

// Next returns a channel which will emit an Event as soon as one of interest occurs
func (self *Server) EtcdWatcher() {
	rootPrefix := "/"
	etcdPath, _ := config.GetString("etcd:path")

	c := etcd.NewClient(etcdPath + rootPrefix)

	log.Info(" [x] Etcd client %s", etcdPath, rootPrefix)

	dir, _ := config.GetString("etcd:directory")
	log.Info(" [x] Etcd Directory %s", dir)

	_, err := c.CreateDir(dir)

	if err != nil {
		log.Error(err)
	}

	etreschan := make(chan *etcd.Response, 1)
	stop := make(chan bool, 1)

	go func() {
		for {
			select {
			case <-stop:
				return
			default:
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
					return
				}
				log.Info(" [x] Sleep-Watch (%s)", rootPrefix+dir)

				time.Sleep(time.Second)

				log.Info(" [x] Slept-Watch (%s)", rootPrefix+dir)
			}
		}

	}()

	go receiverEtcd(etreschan, stop)

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

	stop <- true
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
