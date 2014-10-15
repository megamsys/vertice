package http

import (
	log "code.google.com/p/log4go"
	"github.com/bmizerany/pat"
	"github.com/tsuru/config"
	"net"
	libhttp "net/http"
	"strconv"
	"strings"
	"time"
)

type TimePrecision int
type HttpServer struct {
	conn     net.Listener
	HttpPort int
	//	adminAssetsDir string
	shutdown    chan bool
	readTimeout time.Duration
	p           *pat.PatternServeMux
}

func NewHttpServer() *HttpServer {
	//apiReadTimeout, _ := config.GetString("read-timeout")
	apiHttpPortString, _ := config.GetInt("admin:port")
	self := &HttpServer{}
	self.HttpPort = apiHttpPortString
	//self.adminAssetsDir = config.AdminAssetsDir
	self.shutdown = make(chan bool, 2)
	self.p = pat.New()
	self.readTimeout = 10 * time.Second
	return self
}

func (self *HttpServer) ListenAndServe() {
	var err error
	if self.HttpPort > 0 {
		self.conn, err = net.Listen("tcp", ":"+strconv.Itoa(self.HttpPort))
		if err != nil {
			log.Error("Listen: ", err)
		}
	}
	self.Serve(self.conn)
}

func (self *HttpServer) registerEndpoint(method string, pattern string, f libhttp.HandlerFunc) {
	version, _ := config.GetString("version")
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

func (self *HttpServer) sendCrossOriginHeader(w libhttp.ResponseWriter, r *libhttp.Request) {
	w.WriteHeader(libhttp.StatusOK)
}

func isPretty(r *libhttp.Request) bool {
	return r.URL.Query().Get("pretty") == "true"
}

func (self *HttpServer) query(w libhttp.ResponseWriter, r *libhttp.Request) {

}
