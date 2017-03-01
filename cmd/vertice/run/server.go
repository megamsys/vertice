package run

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	log "github.com/Sirupsen/logrus"
	pp "github.com/megamsys/libgo/cmd"
	"github.com/megamsys/vertice/meta"
	"github.com/megamsys/vertice/subd/deployd"
	"github.com/megamsys/vertice/subd/dns"
	"github.com/megamsys/vertice/subd/docker"
	"github.com/megamsys/vertice/subd/eventsd"
	"github.com/megamsys/vertice/subd/httpd"
	"github.com/megamsys/vertice/subd/marketplacesd"
	"github.com/megamsys/vertice/subd/metricsd"
	"github.com/megamsys/vertice/subd/rancher"
)

// Server represents a container for the metadata and storage data and services.
// It is built using a config and it manages the startup and shutdown of all
// services in the proper order.
type Server struct {
	version string // Build version

	err      chan error
	closing  chan struct{}
	Services []Service

	// Profiling
	CPUProfile string
	MemProfile string
}

// NewServer returns a new instance of Server built from a config.
func NewServer(c *Config, version string) (*Server, error) {
	s := &Server{
		version: version,
		err:     make(chan error),
		closing: make(chan struct{}),
	}

	s.appendDeploydService(c.Meta, c.Deployd)
	s.appendHTTPDService(c.HTTPD)
	s.appendDockerService(c.Meta, c.Docker)
	s.appendMetricsdService(c)
	s.appendEventsdService(c.Meta, c.Events, c.Deployd)
	s.appendRancherService(c.Meta, c.Rancher)
	s.appendMarketplacesService(c.Meta, c.MarketPlaces, c.Deployd)
	s.selfieDNS(c.DNS)
	c.Meta.MkGlobal() //a setter for global meta config
	return s, nil
}

func (s *Server) appendDeploydService(c *meta.Config, d *deployd.Config) {
	e := *d
	if !e.One.Enabled {
		log.Warn("skip oned service.")
		return
	}
	srv := deployd.NewService(c, d)
	s.Services = append(s.Services, srv)
}

func (s *Server) appendHTTPDService(c *httpd.Config) {
	e := *c
	if !e.Enabled {
		log.Warn("skip httpd service.")
		return
	}
	srv := httpd.NewService(c)
	s.Services = append(s.Services, srv)
}

func (s *Server) appendDockerService(c *meta.Config, d *docker.Config) {
	e := *d
	if !e.Docker.Enabled {
		log.Warn("skip dockerd service.")
		return
	}
	srv := docker.NewService(c, d)
	s.Services = append(s.Services, srv)
}

func (s *Server) appendMarketplacesService(c *meta.Config, d *marketplacesd.Config, one *deployd.Config) {
	e := *d
	if !e.Enabled {
		log.Warn("skip marketplacesd service.")
		return
	}
	srv := marketplacesd.NewService(c, d, one)
	s.Services = append(s.Services, srv)
}

func (s *Server) appendMetricsdService(c *Config) {
	if !c.Metrics.Enabled {
		log.Warn("skip metricsd service.")
		return
	}
	srv := metricsd.NewService(c.Meta, c.Deployd, c.Docker, c.Metrics, c.Storage)
	s.Services = append(s.Services, srv)
}

func (s *Server) appendEventsdService(c *meta.Config, e *eventsd.Config, o *deployd.Config) {
	if !e.Enabled {
		log.Warn("skip eventsd service.")
		return
	}
	srv := eventsd.NewService(c, e, o)
	s.Services = append(s.Services, srv)
}

func (s *Server) appendRancherService(c *meta.Config, d *rancher.Config) {
	e := *d
	if !e.Rancher.Enabled {
		log.Warn("skip Rancher service.")
		return
	}
	srv := rancher.NewService(c, d)
	s.Services = append(s.Services, srv)
}

//we are just making the DNS config global
func (s *Server) selfieDNS(c *dns.Config) {
	c.MkGlobal()
}

// Err returns an error channel that multiplexes all out of band errors received from all services.
func (s *Server) Err() <-chan error { return s.err }

// Open opens the meta and data store and all services.
func (s *Server) Open() error {
	if err := func() error {
		//Start profiling, if set.
		startProfile(s.CPUProfile, s.MemProfile)
		//go s.monitorErrorChan(s.?.Err())
		for _, service := range s.Services {
			if err := service.Open(); err != nil {
				return fmt.Errorf("open service: %s", err)
			}
		}
		log.Debug(pp.Colorfy("ō͡≡o˞̶  engine up", "green", "", "bold"))
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

	/*if s.eventHander !=nil {
		s.CloseEventChannel
	}*/
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
