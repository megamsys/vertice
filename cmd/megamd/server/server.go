package server

import (
	"github.com/megamsys/megamd/api/http"
	"github.com/megamsys/megamd/cmd/megamd/server/queue"
	log "code.google.com/p/log4go"
)

type Server struct {
	HttpApi      *http.HttpServer
	QueueServers []*queue.QueueServer
	//AdminServer  *admin.HttpServer
	stopped      bool
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
	var queueInput [2]string
	queueInput[0] = "cloudstandup"
	queueInput[1] = "Events"
	
	// Queue input
	for i := range queueInput {
		listenQueue := queueInput[i]
		queueserver := queue.NewServer(listenQueue)
		go queueserver.ListenAndServe()
	}

	self.HttpApi.ListenAndServe()
	return nil
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
