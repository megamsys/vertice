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
	"github.com/megamsys/libgo/amqp"
	"github.com/megamsys/libgo/db"
	"github.com/megamsys/megamd/api/http"
	"github.com/megamsys/megamd/cmd/megamd/server/queue"
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
