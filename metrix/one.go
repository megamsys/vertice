package metrix

import (
	"encoding/xml"
	"github.com/megamsys/opennebula-go/metrics"
	"github.com/megamsys/vertice/carton"
	"io/ioutil"
	"strconv"
	"time"
)

const (
	OPENNEBULA  = "one"
)

type OpenNebula struct {
	Url          string
	DefaultUnits map[string]string
	RawStatus    []byte
}

func (on *OpenNebula) Prefix() string {
	return "one"
}

func (on *OpenNebula) DeductBill(c *MetricsCollection) (e error) {
	for _, mc := range c.Sensors {
			mkBalance(mc, on.DefaultUnits)
	}
	return
}

func (on *OpenNebula) Collect(c *MetricsCollection) (e error) {
	b, e := on.ReadStatus()
	if e != nil {
		return
	}

	s, e := on.ParseStatus(b)
	if e != nil {
		return
	}
	on.CollectMetricsFromStats(c, s)
	e = on.DeductBill(c)
	return
}

func (on *OpenNebula) ReadStatus() (b []byte, e error) {
	if len(on.RawStatus) == 0 {
		var res []interface{}
		res, e = carton.ProvisionerMap[on.Prefix()].MetricEnvs(time.Now().Add(-10*time.Minute).Unix(), time.Now().Unix(), on.Url, ioutil.Discard)
		if e != nil {
			return
		}
		on.RawStatus = []byte(res[0].(string))
	}

	b = on.RawStatus
	return
}

func (on *OpenNebula) ParseStatus(b []byte) (ons *metrics.OpenNebulaStatus, e error) {
	ons = &metrics.OpenNebulaStatus{}
	e = xml.Unmarshal(b, ons)
	if e != nil {
		return nil, e
	}
	return ons, nil
}

//actually the NewSensor can create trypes based on the event type.
func (on *OpenNebula) CollectMetricsFromStats(mc *MetricsCollection, s *metrics.OpenNebulaStatus) {
	for _, h := range s.History_Records {
		sc := NewSensor(ONE_VM_SENSOR)
		sc.AccountId = h.AccountsId()
		sc.System = on.Prefix()
		sc.Node = h.HostName
		sc.AssemblyId = h.AssemblyId()
		sc.AssemblyName = h.AssemblyName()
		sc.AssembliesId = h.AssembliesId()
		sc.Source = on.Prefix()
		sc.Message = "vm billing"
		sc.Status = h.State()
		sc.AuditPeriodBeginning = time.Unix(h.PStime, 0).String()
		sc.AuditPeriodEnding = time.Unix(h.PEtime, 0).String()
		sc.AuditPeriodDelta = h.Elapsed()
		sc.addMetric(CPU_COST, h.CpuCost(), h.Cpu(), "delta")
		sc.addMetric(MEMORY_COST, h.MemoryCost(), h.Memory(), "delta")
		sc.addMetric(DISK_COST, h.DiskCost(),strconv.FormatInt(h.DiskSize(),10), "delta")
		sc.CreatedAt = time.Now()
		if sc.isBillable() {
			mc.Add(sc)
		}

	}
	return
}
