package metricsd

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/vertice/meta"
	"github.com/megamsys/vertice/metrix"
	"github.com/megamsys/vertice/subd/deployd"
	"github.com/megamsys/vertice/subd/docker"
)

const (
	RIAK = "riak"
)

var OUTPUTS = []string{RIAK}

// Service manages the listener and handler for an HTTP endpoint.
type Service struct {
	err     chan error
	Handler *Handler
	stop    chan struct{}
	Meta    *meta.Config
	Deployd *deployd.Config
	Dockerd *docker.Config
	Config  *Config
}

// NewService returns a new instance of Service.
func NewService(c *meta.Config, one *deployd.Config,doc *docker.Config, f *Config) *Service {
	s := &Service{
		err:     make(chan error),
		Meta:    c,
		Deployd: one,
		Dockerd: doc,
		Config:  f,
	}
	s.Handler = NewHandler()
	return s
}

// Open starts the service
func (s *Service) Open() error {
	log.Info("starting metricsd service")
	if s.stop != nil {
		return nil
	}

	s.stop = make(chan struct{})
	go s.backgroundLoop()
	return nil
}

func (s *Service) backgroundLoop() {
	for {
		select {
		case <-s.stop:
			log.Info("metricsd terminating")
			break
		case <-time.After(time.Duration(s.Config.CollectInterval)):
			s.runMetricsCollectors()
		}
	}

}

func (s *Service) runMetricsCollectors() error {
	output := &metrix.OutputHandler{
		ScyllaAddress: s.Meta.Scylla,
	}

	for _, region := range s.Deployd.One.Regions {
		collectors := map[string]metrix.MetricCollector{
			metrix.OPENNEBULA: &metrix.OpenNebula{Url: region.OneEndPoint},
		}

		mh := &metrix.MetricHandler{}

		for _, collector := range collectors {
			go s.Handler.processCollector(mh, output, collector)
		}
	}

	for _, region := range s.Dockerd.Docker.Regions {
		collectors := map[string]metrix.MetricCollector{
			metrix.DOCKER: &metrix.Swarm{Url: region.SwarmEndPoint},
		}

		mh := &metrix.MetricHandler{}

		for _, collector := range collectors {
			go s.Handler.processCollector(mh, output, collector)
		}
	}



	return nil
}

func (s *Service) Close() error {
	if s.stop == nil {
		return nil
	}
	close(s.stop)
	s.stop = nil
	return nil
}

// Err returns a channel for fatal errors that occur on the listener.
func (s *Service) Err() <-chan error { return s.err }
