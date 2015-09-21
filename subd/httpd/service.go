package httpd

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/megamsys/megamd/api"
	"github.com/megamsys/megamd/subd/httpd/shutdown"
	"gopkg.in/tylerb/graceful.v1"
)

// Service manages the listener and handler for an HTTP endpoint.
type Service struct {
	ln           *graceful.Server
	addr         string
	tls          bool
	certFile     string
	keyFile      string
	err          chan error
	shutdownChan chan bool
	hlr          *negroni.Negroni
}

// NewService returns a new instance of Service.
func NewService(c *Config) *Service {
	s := &Service{
		addr:     c.BindAddress,
		tls:      c.UseTls,
		certFile: c.CertFile,
		keyFile:  c.KeyFile,
		err:      make(chan error),
		hlr:      api.NewNegHandler(),
	}
	return s
}

// Open starts the service
func (s *Service) Open() error {
	log.Infof("Starting httpd service")
	shutdownChan := make(chan bool)
	shutdownTimeout := 10 * 60
	idleTracker := newIdleTracker()
	shutdown.Register(idleTracker)
	shutdown.Register(&api.LogTracker)
	readTimeout := 10 * 60
	writeTimeout := 10 * 60
	srv := &graceful.Server{
		Timeout: time.Duration(shutdownTimeout) * time.Second,
		Server: &http.Server{
			ReadTimeout:  time.Duration(readTimeout) * time.Second,
			WriteTimeout: time.Duration(writeTimeout) * time.Second,
			Addr:         s.addr,
			Handler:      s.hlr,
		},
		ConnState: func(conn net.Conn, state http.ConnState) {
			idleTracker.trackConn(conn, state)
		},
		ShutdownInitiated: func() {
			fmt.Println("megamd-httpd is shutting down, waiting for pending connections to finish.")
			handlers := shutdown.All()
			wg := sync.WaitGroup{}
			for _, h := range handlers {
				wg.Add(1)
				go func(h shutdown.Shutdownable) {
					defer wg.Done()
					fmt.Printf("running shutdown handler for %v...\n", h)
					h.Shutdown()
					fmt.Printf("running shutdown handler for %v. DONE.\n", h)
				}(h)
			}
			wg.Wait()
			close(shutdownChan)
		},
	}

	s.ln = srv
	s.shutdownChan = shutdownChan
	go s.serve()
	return nil
}

// Close closes the underlying listener.
func (s *Service) Close() error {
	gsrv := s.ln
	if s.ln != nil {
		//this is not  graceful stop. and swallows the exception.
		gsrv.Stop(time.Duration(1*60) * time.Second)
	}
	return nil
}

// Err returns a channel for fatal errors that occur on the listener.
func (s *Service) Err() <-chan error { return s.err }

// serve serves the handler from the listener.
func (s *Service) serve() {
	var err error
	if s.tls {
		fmt.Printf("HTTP/TLS server listening at %s...\n", s.addr)
		err = s.ln.ListenAndServeTLS(s.certFile, s.keyFile)
	} else {
		fmt.Printf("HTTP server listening at %s...\n", s.addr)
		err = s.ln.ListenAndServe()
	}
	if err != nil && !strings.Contains(err.Error(), "closed") {
		s.err <- fmt.Errorf("listener failed: addr=%s, err=%s", s.addr, err)
	}
	<-s.shutdownChan
}
