package metrix

import (
	"fmt"
	"github.com/megamsys/vertice/carton"
	"io/ioutil"
	"time"
)

const DOCKER = "docker"

type Swarm struct {
	Url       string
	DefaultUnits map[string]string
	RawStatus []interface{}
}

type Stats struct {
	ContainerId  string
	MemoryUsage  uint64 //in bytes
	SystemMemory uint64
	CPUStats     CPUStats //in percentage of total cpu used
	PreCPUStats  CPUStats
	NetworkIn    uint64
	NetworkOut   uint64
	AccountId    string
	AssemblyId   string
	AssemblyName string
	AssembliesId string
	Status       string
	AuditPeriod  time.Time
}

type CPUStats struct {
	PercpuUsage       []uint64
	UsageInUsermode   uint64
	TotalUsage        uint64
	UsageInKernelmode uint64
	SystemCPUUsage    uint64
}

func (s *Swarm) Prefix() string {
	return "docker"
}

func (s *Swarm) Collect(c *MetricsCollection) (e error) {
	fmt.Println(s)
	e = s.ReadStatus()
	if e != nil {
		return
	}

	stats, e := s.ParseStatus(s.RawStatus)
	if e != nil {
		return
	}
	s.CollectMetricsFromStats(c, stats)

	e = s.DeductBill(c)
	return
}

func (s *Swarm) DeductBill(c *MetricsCollection) (e error) {
	for _, mc := range c.Sensors {
		if mc.AccountId != "" && mc.AssemblyName != "" {
			mkBalance(mc, s.DefaultUnits)
		}
	}
	return
}

func (s *Swarm) ParseStatus(a []interface{}) ([]*Stats, error) {
	var stats []*Stats
	for _, v := range a {
		f, ok := v.(*Stats)
		if !ok {
			fmt.Println("failed to converter")
		}
		stats = append(stats, f)
	}
	return stats, nil
}

func (s *Swarm) ReadStatus() (e error) {
	s.RawStatus, e = carton.ProvisionerMap[s.Prefix()].MetricEnvs(time.Now().Add(-10*time.Minute).Unix(), time.Now().Unix(), s.Url, ioutil.Discard)
	if e != nil {
		return
	}
	return
}

//actually the NewSensor can create trypes based on the event type.
func (s *Swarm) CollectMetricsFromStats(mc *MetricsCollection, stats []*Stats) {
	for _, h := range stats {
		sc := NewSensor("compute.container.exists")
		sc.AccountId = h.AccountId
		sc.System = s.Prefix()
		sc.Node = ""
		sc.AssemblyId = h.AssemblyId
		sc.AssemblyName = h.AssemblyName
		sc.AssembliesId = h.AssembliesId
		sc.Source = s.Prefix()
		sc.Message = "container billing"
		sc.Status = h.Status
		sc.AuditPeriodBeginning = time.Now().Add(-10 * time.Minute).String()
		sc.AuditPeriodEnding = time.Now().String()
		sc.AuditPeriodDelta = time.Now().String()
		//have calculate the cpu used percentage from 	CPUStats  PreCPUStats
		sc.addMetric(CPU_COST, "2", "0.021", "delta")
		sc.addMetric(MEMORY_COST, "2", "450", "delta")
		mc.Add(sc)
	}
	return
}
