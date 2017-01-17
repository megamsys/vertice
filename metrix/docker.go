package metrix

import (
	"fmt"
	"github.com/megamsys/vertice/carton"
	"io/ioutil"
	"strconv"
	"time"
)

const DOCKER = "docker"

type Swarm struct {
	Url            string
	DefaultUnits map[string]string
	RawStatus      []interface{}
}

type Stats struct {
	ContainerId    string
	Image          string
	MemoryUsage    uint64 //in bytes
	CPUUnitCost    string
	MemoryUnitCost string
	SystemMemory   uint64
	CPUStats       CPUStats //in percentage of total cpu used
	PreCPUStats    CPUStats
	NetworkIn      uint64
	NetworkOut     uint64
	AccountId      string
	AssemblyId     string
	QuotaId        string
	AssemblyName   string
	AssembliesId   string
	Status         string
	AuditPeriod    time.Time
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
		if mc.AccountId != "" && mc.AssemblyId != "" {
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
	s.RawStatus, e = carton.ProvisionerMap[s.Prefix()].MetricEnvs(time.Now().Add(-MetricsInterval).Unix(), time.Now().Unix(), s.Url, ioutil.Discard)
	if e != nil {
		return
	}
	return
}

//actually the NewSensor can create types based on the event type.
func (s *Swarm) CollectMetricsFromStats(mc *MetricsCollection, stats []*Stats) {

	for _, h := range stats {
		if !(len(h.QuotaId) > 0) {
			cpuDelta := float64((float64(h.CPUStats.TotalUsage) - float64(h.PreCPUStats.TotalUsage)))
			systemDelta := float64((float64(h.CPUStats.SystemCPUUsage) - float64(h.PreCPUStats.SystemCPUUsage)))
			cpu_usage := (cpuDelta / systemDelta) * float64(len(h.CPUStats.PercpuUsage)) * 100.0
			sc := NewSensor(DOCKER_CONTAINER_SENSOR)
			sc.AccountId = h.AccountId
			sc.System = s.Prefix()
			sc.Node = ""
			sc.AssemblyId = h.AssemblyId
			sc.AssemblyName = h.AssemblyName
			sc.AssembliesId = h.AssembliesId
			sc.Source = s.Prefix()
			sc.Message = "container billing"
			sc.Status = h.Status
			sc.AuditPeriodBeginning = time.Now().Add(-MetricsInterval).Format(time.RFC3339)
			sc.AuditPeriodEnding = time.Now().Format(time.RFC3339)
			sc.AuditPeriodDelta = MetricsInterval.String()
			//have calculate the cpu used percentage from 	CPUStats  PreCPUStats
			sc.addMetric(CPU_COST, h.CPUUnitCost, strconv.FormatFloat(cpu_usage, 'f', 6, 64), "delta")
			sc.addMetric(MEMORY_COST, h.MemoryUnitCost, strconv.FormatFloat(float64(h.MemoryUsage/1024.0/1024.0), 'f', 6, 64), "delta")
			mc.Add(sc)
			sc.CreatedAt = time.Now()
			if sc.isBillable() {
				mc.Add(sc)
			}
		}

	}
	return
}
