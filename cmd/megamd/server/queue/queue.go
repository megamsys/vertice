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
package queue

import (
	"github.com/megamsys/libgo/amqp"
	"github.com/megamsys/megamd/global"
	"github.com/megamsys/megamd/coordinator"
)

type QueueServer struct {
	ListenAddress string
	chann         chan []byte
	shutdown      chan bool
}

//interface arguments
func NewServer(listenAddress string) *QueueServer {
	global.LOG.Info("Create New queue server")
	self := &QueueServer{}

	self.ListenAddress = listenAddress
	self.shutdown = make(chan bool, 1)
	global.LOG.Info("queue: %s", self)
	return self
}

func (self *QueueServer) ListenAndServe() {
	factor, err := amqp.Factory()
	if err != nil {
		global.LOG.Error("Failed to get the queue instance: %s", err)
	}

	pubsub, err := factor.Get(self.ListenAddress)
	if err != nil {
		global.LOG.Error("Failed to get the queue instance: %s", err)
	}

	msgChan, _ := pubsub.Sub()
	for msg := range msgChan {
		global.LOG.Info(" [x] %q", msg)
		global.LOG.Info("I am in! ")
		coordinator.NewCoordinator(msg, self.ListenAddress)
	}
	
	global.LOG.Info("Handling message %v", msgChan)
	self.chann = msgChan

	//self.Serve()
}
