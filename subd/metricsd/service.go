package metricsd

import (
	log "github.com/Sirupsen/logrus"
	constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/vertice/meta"
	"github.com/megamsys/vertice/metrix"
	"github.com/megamsys/vertice/storage"
	"github.com/megamsys/vertice/subd/deployd"
	"github.com/megamsys/vertice/subd/docker"
	"strconv"
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
}

// NewService returns a new instance of Service.
func NewService(c *meta.Config, one *deployd.Config, doc *docker.Config, f *Config, strg *storage.Config) *Service {
	s := &Service{
		err:     make(chan error),
		Meta:    c,
		Deployd: one,
		Dockerd: doc,
		Config:  f,
		Storage: strg,
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
	skews := make(map[string]string, 0)
	skews[constants.ENABLED] = strconv.FormatBool(s.Config.Skews.Enabled)
	skews[constants.SOFT_LIMIT] = s.Config.Skews.SoftLimit
	skews[constants.SOFT_GRACEPERIOD] = s.Config.Skews.SoftGracePeriod.String()
	skews[constants.HARD_LIMIT] = s.Config.Skews.HardLimit
	skews[constants.HARD_GRACEPERIOD] = s.Config.Skews.HardGracePeriod.String()
	metrix.MetricsInterval = time.Duration(s.Config.CollectInterval)

	if s.Deployd.One.Enabled {
		s.onedCollectors(output, skews)
	}

	if s.Dockerd.Docker.Enabled {
		s.dockerCollectors(output, skews)
	}

	if s.Storage.Enabled {
		s.storageCollectors(output, skews)
	}

	if s.Config.Backups.Enabled {
		s.backupsCollectors(output, skews)
	}

	if s.Config.Snapshots.Enabled {
		s.snapshotsCollectors(output, skews)
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

func (s *Service) onedCollectors(output *metrix.OutputHandler, skews map[string]string) {
	// One VirtualMachine Metrics collectors
	if s.Deployd.One.Enabled {
		for _, region := range s.Deployd.One.Regions {
			collectors := map[string]metrix.MetricCollector{
				metrix.OPENNEBULA: &metrix.OpenNebula{
					Url:          region.OneEndPoint,
					Region:       region.OneZone,
					DefaultUnits: map[string]string{metrix.MEMORY_UNIT: region.MemoryUnit, metrix.CPU_UNIT: region.CpuUnit, metrix.DISK_UNIT: region.DiskUnit},
					SkewsActions: skews,
				},
			}

			mh := &metrix.MetricHandler{}

			for _, collector := range collectors {
				go s.Handler.processCollector(mh, output, collector)
			}
		}
	}

}

func (s *Service) dockerCollectors(output *metrix.OutputHandler, skews map[string]string) {
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

func (s *Service) storageCollectors(output *metrix.OutputHandler, skews map[string]string) {
	if s.Storage.RgwStorage.Enabled {
		// Ceph RadosGW (storage buckets) Metrics collectors
		for _, region := range s.Storage.RgwStorage.Regions {
			collectors := map[string]metrix.MetricCollector{
				metrix.CEPHRGW: &metrix.CephRGWStats{Url: region.EndPoint,
					DefaultUnits: map[string]string{metrix.STORAGE_UNIT: region.StorageUnit, metrix.STORAGE_COST_PER_HOUR: region.CostPerHour},
					AdminUser:    region.AdminUser,
					MasterKey:    s.Meta.MasterKey,
					AccessKey:    region.AdminAccess,
					SecretKey:    region.AdminSecret,
				},
			}

			mh := &metrix.MetricHandler{}

			for _, collector := range collectors {
				go s.Handler.processCollector(mh, output, collector)
			}

		}
	}
}

func (s *Service) snapshotsCollectors(output *metrix.OutputHandler, skews map[string]string) {
	// snapshots collectors
	collectors := map[string]metrix.MetricCollector{
		metrix.SNAPSHOTS: &metrix.Snapshots{
			DefaultUnits: map[string]string{metrix.STORAGE_UNIT: s.Config.Snapshots.StorageUnit, metrix.STORAGE_COST_PER_HOUR: s.Config.Snapshots.CostPerHour},
		},
	}
	mh := &metrix.MetricHandler{}

	for _, collector := range collectors {
		go s.Handler.processCollector(mh, output, collector)
	}
}

func (s *Service) backupsCollectors(output *metrix.OutputHandler, skews map[string]string) {
	// snapshots collectors
	collectors := map[string]metrix.MetricCollector{
		metrix.BACKUPS: &metrix.Backups{
			DefaultUnits: map[string]string{metrix.STORAGE_UNIT: s.Config.Backups.StorageUnit, metrix.STORAGE_COST_PER_HOUR: s.Config.Backups.CostPerHour},
		},
	}
	mh := &metrix.MetricHandler{}

	for _, collector := range collectors {
		go s.Handler.processCollector(mh, output, collector)
	}
}
