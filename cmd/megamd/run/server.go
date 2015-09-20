package run

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/megamd/meta"
	"github.com/megamsys/megamd/subd/deployd"
	"github.com/megamsys/megamd/subd/docker"
	"github.com/megamsys/megamd/subd/httpd"
)

// Server represents a container for the metadata and storage data and services.
// It is built using a Config and it manages the startup and shutdown of all
// services in the proper order.
type Server struct {
	version string // Build version

	err     chan error
	closing chan struct{}

	Hostname    string
	BindAddress string

	Services []Service

	// Profiling
	CPUProfile string
	MemProfile string
}

// NewServer returns a new instance of Server built from a config.
func NewServer(c *Config, version string) (*Server, error) {
	// Construct base meta store and data store.
	s := &Server{
		version: version,
		err:     make(chan error),
		closing: make(chan struct{}),

		Hostname:    c.Meta.Hostname,
		BindAddress: c.Meta.BindAddress,
	}

	// Append services.
	s.appendDeploydService(c.Meta, c.Deployd)
	s.appendHTTPDService(c.HTTPD)
	s.appendDockerService(c.Meta, c.Docker)
	//s.appendEventsTransporter(c.Meta)
	return s, nil
}

func (s *Server) appendDeploydService(c *meta.Config, d *deployd.Config) {
	srv := deployd.NewService(c, d)
	//	srv.ProvisioningWriter = s.ProvisioningWriter
	s.Services = append(s.Services, srv)
}

func (s *Server) appendDockerService(c *meta.Config, d *docker.Config) {
	if !d.Enabled {
		log.Warn("skip docker service.")
		return
	}
	srv := docker.NewService(c, d)
	//	srv.SwarmExecutor = s.SwarmExecutor
	//	srv.ProvisioningWriter = s.ProvisioningWriter
	s.Services = append(s.Services, srv)
}

func (s *Server) appendHTTPDService(c *httpd.Config) {
	e := *c
	if !e.Enabled {
		log.Warn("skip httpd service.")
		return
	}
	srv := httpd.NewService(c)
	//	srv.Handler.QueryExecutor = s.QueryExecutor

	s.Services = append(s.Services, srv)
}

// Err returns an error channel that multiplexes all out of band errors received from all services.
func (s *Server) Err() <-chan error { return s.err }

// Open opens the meta and data store and all services.
func (s *Server) Open() error {
	if err := func() error {
		// Start profiling, if set.
		startProfile(s.CPUProfile, s.MemProfile)

		/*	host, port, err := s.hostAddr()
			if err != nil {
				return err
			}
		*/
		//		go s.monitorErrorChan(s.?.Err())

		for _, service := range s.Services {
			if err := service.Open(); err != nil {
				return fmt.Errorf("open service: %s", err)
			}
		}
		log.Debug("services started")

		return nil

	}(); err != nil {
		s.Close()
		return err
	}

	return nil
}

// Close shuts down the meta and data stores and all services.
func (s *Server) Close() error {
	stopProfile()

	for _, service := range s.Services {
		service.Close()
	}

	if s.closing != nil {
		close(s.closing)
	}
	return nil
}

// monitorErrorChan reads an error channel and resends it through the server.
func (s *Server) monitorErrorChan(ch <-chan error) {
	for {
		select {
		case err, ok := <-ch:
			if !ok {
				return
			}
			s.err <- err
		case <-s.closing:
			return
		}
	}
}

// hostAddr returns the host and port that remote nodes will use to reach this
// node.
func (s *Server) hostAddr() (string, string, error) {
	// Resolve host to address.
	_, port, err := net.SplitHostPort(s.BindAddress)
	if err != nil {
		return "", "", fmt.Errorf("split bind address: %s", err)
	}

	host := s.Hostname

	// See if we might have a port that will override the BindAddress port
	if host != "" && host[len(host)-1] >= '0' && host[len(host)-1] <= '9' && strings.Contains(host, ":") {
		hostArg, portArg, err := net.SplitHostPort(s.Hostname)
		if err != nil {
			return "", "", err
		}

		if hostArg != "" {
			host = hostArg
		}

		if portArg != "" {
			port = portArg
		}
	}
	return host, port, nil
}

// Service represents a service attached to the server.
type Service interface {
	Open() error
	Close() error
}

// prof stores the file locations of active profiles.
var prof struct {
	cpu *os.File
	mem *os.File
}

// StartProfile initializes the cpu and memory profile, if specified.
func startProfile(cpuprofile, memprofile string) {
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Errorf("cpuprofile: %v", err)
		}
		log.Infof("writing CPU profile to: %s", cpuprofile)
		prof.cpu = f
		pprof.StartCPUProfile(prof.cpu)
	}

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			log.Errorf("memprofile: %v", err)
		}
		log.Infof("writing mem profile to: %s", memprofile)
		prof.mem = f
		runtime.MemProfileRate = 4096
	}

}

// StopProfile closes the cpu and memory profiles if they are running.
func stopProfile() {
	if prof.cpu != nil {
		pprof.StopCPUProfile()
		prof.cpu.Close()
		log.Infof("CPU profile stopped")
	}
	if prof.mem != nil {
		pprof.Lookup("heap").WriteTo(prof.mem, 0)
		prof.mem.Close()
		log.Infof("mem profile stopped")
	}
}
