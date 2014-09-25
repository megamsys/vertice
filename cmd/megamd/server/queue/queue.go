package queue

import (
	"github.com/megamsys/libgo/amqp"
	log "code.google.com/p/log4go"
)

type QueueServer struct {
	listenAddress string
	chann          chan []byte
	shutdown      chan bool
}
//interface arguments
func NewServer(listenAddress string) *QueueServer {
	log.Info("Create New queue server")
	self := &QueueServer{}

	self.listenAddress = listenAddress
	self.shutdown = make(chan bool, 1)
     log.Info(self)
	return self
}


func (self *QueueServer) ListenAndServe() {
	
	log.Info("------------------------------------------")
    log.Info("Starting Queue listen")
	factor, _ := amqp.Factory()
	pubsub, _ := factor.Get(self.listenAddress)
	msgChan, _ := pubsub.Sub()
	self.chann = msgChan
	
}

