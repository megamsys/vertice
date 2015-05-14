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
package server

import (
	log "code.google.com/p/log4go"
	"encoding/json"
	"github.com/megamsys/libgo/amqp"
	"github.com/megamsys/libgo/db"
	"github.com/megamsys/libgo/etcd"
	"github.com/megamsys/megamd/api/http"
	"github.com/megamsys/megamd/global"
	"github.com/megamsys/megamd/cmd/megamd/server/queue"
	"github.com/tsuru/config"
	"time"
	"os"
	"fmt"
	"net"
	"net/url"
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
	var queueInput [3]string
	queueInput[0] = "cloudstandup"
	queueInput[1] = "events"
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

func (self *Server) EtcdWatcher() {
	rootPrefix := "/"
	etcdPath, _ := config.GetString("etcd:path")

	c := etcd.NewClient([]string{etcdPath})
	success := c.SyncCluster()
	if !success {
		log.Debug("cannot sync machines")
	}

	for _, m := range c.GetCluster() {
		u, err := url.Parse(m)
		if err != nil {
			log.Debug(err)
		}
		if u.Scheme != "http" {
			log.Debug("scheme must be http")
		}

		host, _, err := net.SplitHostPort(u.Host)
		if err != nil {
			log.Debug(err)
		}
		if host != "127.0.0.1" {
			log.Debug("Host must be 127.0.0.1")
		}
	}
	
	etcdNetworkPath, _ := config.GetString("etcd:networkpath")
	conn, connerr := c.Dial("tcp", etcdNetworkPath)
    log.Debug("client %v", c)
    log.Debug("connection %v", conn)
    log.Debug("connection error %v", connerr)
    
    if conn != nil {
	
	   log.Info(" [x] Etcd client %s", etcdPath, rootPrefix)

	   dir, _ := config.GetString("etcd:directory")
	   log.Info(" [x] Etcd Directory %s", dir)

	   stop := make(chan bool, 0)

	   go func() {
		for {
			select {
			case <-stop:
				return
			default:
			   _, err1 := c.CreateDir(dir)

	           if err1 != nil {
		         //  log.Error(err1)
	           }
				etreschan := make(chan *etcd.Response, 1)
			     go receiverEtcd(etreschan, stop) 
			 	_, err := c.Watch(rootPrefix+dir, 0, true, etreschan, stop)

				if err != nil {
					//log.Info(" [x] Watched Error (%s)", rootPrefix+dir)
				//	log.Error(err)
				//	return
				}
				if err != etcd.ErrWatchStoppedByUser {
				//	log.Error("Watch returned a non-user stop error")
				}
				//log.Info(" [x] Sleep-Watch (%s)", rootPrefix+dir)

				time.Sleep(time.Second)

				//log.Info(" [x] Slept-Watch (%s)", rootPrefix+dir)
			}
		} 

	}()
	
	} else {
  	 fmt.Fprintf(os.Stderr, "Error: %v\n Please start etcd deamon.\n", connerr)
         os.Exit(1)
  }
}

/**
In this goroutine received the message from channel then to export the message to handler, 
and this goroutine is close when the message is nil. 
**/
func receiverEtcd(c chan *etcd.Response, stop chan bool) {
	for {
		select {
		case msg := <-c:
			//log.Info(" [x] Handing etcd master response (%v)", msg)
			if msg != nil {
				handlerEtcd(msg)
			} else {
			//	log.Info(" [x] Nil - Handling etcd master response (%v)", msg)
				return
			}
		}
	}
	stop <- false
}

func handlerEtcd(msg *etcd.Response) {
	log.Info(" [x] Really Handle etcd response (%s)", msg.Node.Key)

	asm := &global.Assembly{}

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
		pair, perr := global.ParseKeyValuePair(asm.Inputs, "domain")
		if perr != nil {
			log.Error("Failed to get the domain value : %s", perr)
		}
		asmname := asm.Name+"."+pair.Value
		//asmname := asm.Name
		publisher(asmname, string(mapB))
	}
}

func publisher(key string, json string) {
	factor, aerr := amqp.Factory()
	if aerr != nil {
		log.Error("Failed to get the queue instance: %s", aerr)
		return
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
