package http

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	libhttp "net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"github.com/tsuru/config"
)

type HttpServer struct {
	conn           net.Listener
	HttpPort       string
//	adminAssetsDir string
	shutdown       chan bool
	readTimeout    time.Duration
	p              *pat.PatternServeMux
}

func NewHttpServer() *HttpServer {
	apiReadTimeout, _ = config.GetString("read-timeout")
	apiHttpPortString, _ = config.GetString("admin:port")
	self := &HttpServer{}
	self.HttpPort = apiHttpPortString
	//self.adminAssetsDir = config.AdminAssetsDir	
	self.shutdown = make(chan bool, 2)
	self.p = pat.New()
	self.readTimeout = apiReadTimeout
	return self
}

func (self *HttpServer) ListenAndServe() {
	var err error
	if self.httpPort != "" {
		self.conn, err = net.Listen("tcp", self.httpPort)
		if err != nil {
			log.Error("Listen: ", err)
		}
	}
	self.Serve(self.conn)
}

func (self *HttpServer) registerEndpoint(method string, pattern string, f libhttp.HandlerFunc) {
	version, _ := config.GetString("version")
	//version := self.clusterConfig.GetLocalConfiguration().Version
	switch method {
	case "get":
		self.p.Get(pattern, CompressionHeaderHandler(f, version))
	case "post":
		self.p.Post(pattern, HeaderHandler(f, version))
	case "del":
		self.p.Del(pattern, HeaderHandler(f, version))
	}
	self.p.Options(pattern, HeaderHandler(self.sendCrossOriginHeader, version))
}

func (self *HttpServer) Serve(listener net.Listener) {
	defer func() { self.shutdown <- true }()

	self.conn = listener

	// Run the given query and return an array of series or a chunked response
	// with each batch of points we get back
	self.registerEndpoint("get", "/index", self.query)

	
	self.serveListener(listener, self.p)
}

func (self *HttpServer) serveListener(listener net.Listener, p *pat.PatternServeMux) {
	srv := &libhttp.Server{Handler: p, ReadTimeout: self.readTimeout}
	if err := srv.Serve(listener); err != nil && !strings.Contains(err.Error(), "closed network") {
		panic(err)
	}
}

