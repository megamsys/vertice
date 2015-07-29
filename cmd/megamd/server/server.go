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
	"encoding/json"
	"fmt"
	"os"

	log "code.google.com/p/log4go"
	"github.com/megamsys/libgo/amqp"
	"github.com/megamsys/libgo/db"
	"github.com/megamsys/megamd/api/http"
	"github.com/megamsys/megamd/cmd/megamd/server/queue"
	"github.com/megamsys/megamd/global"
	"github.com/tsuru/config"
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
	queueInput[2] = "dockerstate"
	self.Checker()
	self.IPInit()

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
	Dial string `json:"dial"`
}

func (self *Server) Checker() {
	log.Info("verifying rabbitmq")
	factor, err := amqp.Factory()
	if err != nil {
		log.Error("Error: %v\nFailed to get the queue", err)
	}

	_, connerr := factor.Dial()
	if connerr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n Please start rabbitmq service.\n", connerr)
		os.Exit(1)
	}
	log.Info("rabbitmq connected [ok]")

	log.Info("verifying riak")

	rconn, rerr := db.Conn("connection")
	if rerr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n Please start Riak service.\n", rerr)
		os.Exit(1)
	}

	data := "sampledata"
	ferr := rconn.StoreObject("sampleobject", data)
	if ferr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n Please start Riak service.\n", ferr)
		os.Exit(1)
	}
	defer rconn.Close()
	log.Info("riak connected [ok]")

}

func (self *Server) IPInit() {

	rconn, rerr := db.Conn("ipindex")
	if rerr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n Please start Riak service.\n", rerr)
		os.Exit(1)
	}

	index := global.IPIndex{}
	subnet, _ := config.GetString("swarm:subnet")
	_, err := index.Get(global.IPINDEXKEY)
	if err != nil {
		data := &global.IPIndex{Ip: subnet, Subnet: subnet, Index: 1}
		res2B, _ := json.Marshal(data)
		ferr := rconn.StoreObject(global.IPINDEXKEY, string(res2B))
		if ferr != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n Please start Riak service.\n", ferr)
			os.Exit(1)
		}
	}

	defer rconn.Close()
}

func (self *Server) Stop() {
	if self.stopped {
		return
	}
	log.Info("Bye. tata.")
	self.stopped = true
	//self.HttpApi.Close()

}
