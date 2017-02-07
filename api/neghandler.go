package api

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/googollee/go-socket.io"
	"github.com/megamsys/libgo/cmd"
	"golang.org/x/net/websocket"
	"github.com/rs/cors"
	"net/http"
)

type MegdHandler struct {
	method string
	path   string
	h      http.Handler
}

var megdHandlerList []MegdHandler

//RegisterHandler inserts a handler on a list of handlers
func RegisterHandler(path string, method string, h http.Handler) {
	var th MegdHandler
	th.path = path
	th.method = method
	th.h = h
	megdHandlerList = append(megdHandlerList, th)
}

// RunServer starts vertice httpd server.
func NewNegHandler() *negroni.Negroni {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	m := &delayedRouter{}
	for _, handler := range megdHandlerList {
		m.Add(handler.method, handler.path, handler.h)
	}

	socketServer, err := socketio.NewServer(nil)
	if err != nil {
		log.Debugf(cmd.Colorfy("  > [socket] ", "red", "", "bold") + fmt.Sprintf("Error starting socket server : %v", err))
	}

	m.Add("Get", "/", Handler(index))
	//m.Add("Get", "/logs", Handler(logs))
	m.Add("Post", "/logs/", socketServer)
	m.Add("Get", "/logs/", socketServer)
	m.Add("Get", "/ping", Handler(ping))
	m.Add("Get", "/vnc/", Handler(vnc))

	socketHandler(socketServer)

	// Shell also doesn't use {app} on purpose. Middlewares don't play well
	// with websocket.
	m.Add("Get", "/shell/{email}/{asmsid}/{id}", websocket.Handler(remoteShellHandler))

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(c)
	n.Use(newLoggerMiddleware())
	n.UseHandler(m)
	n.Use(negroni.HandlerFunc(contextClearerMiddleware))
	n.Use(negroni.HandlerFunc(flushingWriterMiddleware))
	n.Use(negroni.HandlerFunc(errorHandlingMiddleware))
	n.Use(negroni.HandlerFunc(authTokenMiddleware))
	n.UseHandler(http.HandlerFunc(runDelayedHandler))
	return n
}
