package server

import (
	"fmt"
	"log"
	"runtime"
	"time"
    "github.com/tsuru/config"
	"github.com/megamsys/megamd/api/http"
	"github.com/megamsys/megamd/cmd/megamd/server/queue"
	
)

type Server struct {
	HttpApi      *http.HttpServer
	QueueServers []*queue.Server
	//AdminServer  *admin.HttpServer
	stopped      bool
}

func NewServer() (*Server, error) {
	
	newClient := func(connectString string) cluster.ServerConnection {
		protobuf_timeout, _ = config.GetString("protobuf_timeout")
		return coordinator.NewProtobufClient(connectString, protobuf_timeout)
	}

	httpApi := http.NewHttpServer()
//	httpApi.EnableSsl(config.ApiHttpSslPortString(), config.ApiHttpCertPath)

//	adminServer := admin.NewHttpServer(config.AdminAssetsDir, config.AdminHttpPortString())

	return &Server{
		HttpApi:     httpApi,
		}, nil
		//AdminServer: adminServer}, nil

}

func (self *Server) ListenAndServe() error {
	log.Info("Starting admin interface on port %d", self.HttpPort)
	//go self.AdminServer.ListenAndServe()

	// Queue input
	for _, queueInput := range self.QueueServers {
		port := queueInput.Port
		queue := queueInput.Queue

		if port <= 0 || port >= 65536 {
			log.Warn("Cannot start queue server on port %d. please check your configuration", port)
			continue
		} else if queue == "" {
			log.Warn("Cannot start queue server for queue=\"\".  please check your configuration")
			continue
		}

		log.Info("Starting Queue Listener on port %d to queue %s", port, queue)

		addr := port

		server := queue.NewServer(addr)
		//self.UdpServers = append(self.UdpServers, server)
		go server.ListenAndServe()
	}

    log.Info("Starting Http Api server on port %d", self.HttpPort)
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
	self.HttpApi.Close()
	log.Info("Api server stopped")

	//log.Info("Stopping admin server")
	//self.AdminServer.Close()
	//log.Info("admin server stopped")

}
