package metricsd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/vertice/meta"
	"github.com/megamsys/vertice/metrix"
	"github.com/megamsys/vertice/storage"
	"github.com/megamsys/vertice/snapshots"
	"github.com/megamsys/vertice/subd/deployd"
	"github.com/megamsys/vertice/subd/docker"
	"time"

)

// Service manages the listener and handler for an HTTP endpoint.
type Service struct {
	err     chan error
	Handler *Handler
	stop    chan struct{}
	Meta    *meta.Config
	Deployd *deployd.Config
	Dockerd *docker.Config
	Config  *Config
	Storage *storage.Config
	Snapshots *snapshots.Config
}

// NewService returns a new instance of Service.
func NewService(c *meta.Config, one *deployd.Config, doc *docker.Config, f *Config, strg *storage.Config, snp *snapshots.Config) *Service {
	s := &Service{
		err:     make(chan error),
		Meta:    c,
		Deployd: one,
		Dockerd: doc,
		Config:  f,
		Storage: strg,
		Snapshots: snp,
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
		ScyllaAddress: s.Meta.Api,
	}

	metrix.MetricsInterval = time.Duration(s.Config.CollectInterval)

  if s.Deployd.One.Enabled {
		s.onedCollectors(output)
	}

	if s.Dockerd.Docker.Enabled {
		s.dockerCollectors(output)
	}

	if s.Storage.Enabled {
		s.storageCollectors(output)
	}

	if s.Snapshots.Enabled {
		s.snapshotsCollectors(output)
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


func (s *Service) onedCollectors(output *metrix.OutputHandler) {
	// One VirtualMachine Metrics collectors
	 if s.Deployd.One.Enabled {
		 for _, region := range s.Deployd.One.Regions {
			 collectors := map[string]metrix.MetricCollector{
				 metrix.OPENNEBULA: &metrix.OpenNebula{
					 Url: region.OneEndPoint,
					 DefaultUnits: map[string]string{metrix.MEMORY_UNIT: region.MemoryUnit, metrix.CPU_UNIT: region.CpuUnit, metrix.DISK_UNIT: region.DiskUnit},
				 },
			 }

			 mh := &metrix.MetricHandler{}

			 for _, collector := range collectors {
				 go s.Handler.processCollector(mh, output, collector)
			 }
		 }
	 }

}


func (s *Service) dockerCollectors(output *metrix.OutputHandler) {
	// Docker container Metrics collectors
  if s.Dockerd.Docker.Enabled {
 	 for _, region := range s.Dockerd.Docker.Regions {
 		 collectors := map[string]metrix.MetricCollector{
 			 metrix.DOCKER: &metrix.Swarm{Url: region.SwarmEndPoint, DefaultUnits: map[string]string{metrix.MEMORY_UNIT: region.MemoryUnit, metrix.CPU_UNIT: region.CpuUnit, metrix.DISK_UNIT: region.DiskUnit}},
 		 }

 		 mh := &metrix.MetricHandler{}

 		 for _, collector := range collectors {
 			 go s.Handler.processCollector(mh, output, collector)
 		 }

 	 }
  }
}

func (s *Service) storageCollectors(output *metrix.OutputHandler) {
	if s.Storage.RgwStorage.Enabled {
 	// Ceph RadosGW (storage buckets) Metrics collectors
 	for _, region := range s.Storage.RgwStorage.Regions {
 		collectors := map[string]metrix.MetricCollector{
 			metrix.CEPHRGW: &metrix.CephRGWStats{Url: region.EndPoint,
 				DefaultUnits: map[string]string{metrix.STORAGE_UNIT: region.StorageUnit,  metrix.STORAGE_COST_PER_HOUR: region.CostPerHour},
 				AdminUser: region.AdminUser,
 				MasterKey: s.Meta.MasterKey,
 				AccessKey: region.AdminAccess,
 				SecretKey: region.AdminSecret,
 			},
 		}

 		mh := &metrix.MetricHandler{}

 		for _, collector := range collectors {
 			go s.Handler.processCollector(mh, output, collector)
 		}

 	}
  }
}

func (s *Service) snapshotsCollectors(output *metrix.OutputHandler) {
 	// snapshots collectors
 		collectors := map[string]metrix.MetricCollector{
 			metrix.SNAPSHOTS: &metrix.Snapshots{
				DefaultUnits: map[string]string{metrix.STORAGE_UNIT: s.Snapshots.StorageUnit, metrix.STORAGE_COST_PER_HOUR: s.Snapshots.CostPerHour},
 		  },
	  }
 		mh := &metrix.MetricHandler{}

 		for _, collector := range collectors {
 			go s.Handler.processCollector(mh, output, collector)
 		}
}
