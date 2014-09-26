package queue

import (
	"github.com/megamsys/libgo/amqp"
	log "code.google.com/p/log4go"
	"github.com/megamsys/megamd/coordinator"
)

type QueueServer struct {
	ListenAddress string
	chann          chan []byte
	shutdown      chan bool
}
//interface arguments
func NewServer(listenAddress string) *QueueServer {
	log.Info("Create New queue server")
	self := &QueueServer{}

	self.ListenAddress = listenAddress
	self.shutdown = make(chan bool, 1)
     log.Info(self)
	return self
}



func (self *QueueServer) ListenAndServe() {
	factor, err := amqp.Factory()
	if err != nil {
		log.Error("Failed to get the queue instance: %s", err)
	}
	
	pubsub, err := factor.Get(self.ListenAddress)
	if err != nil {
		log.Error("Failed to get the queue instance: %s", err)
	}
	
	msgChan, _ := pubsub.Sub()
	for msg := range msgChan {
			log.Info(" [x] %q", msg)
			coordinator.NewCoordinator(msg, self.ListenAddress)
		}
	log.Info("Handling message %v", msgChan)
	self.chann = msgChan
	
	//self.Serve()
}



