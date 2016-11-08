package metrix

import (
	"encoding/xml"
	"github.com/megamsys/opennebula-go/metrics"
	"github.com/megamsys/vertice/carton"
	"io/ioutil"
	"time"
)

const OPENNEBULA = "one"

type OpenNebula struct {
	Url       string
	RawStatus []byte
}

func (on *OpenNebula) Prefix() string {
	return "one"
}

func (on *OpenNebula) DeductBill(c *MetricsCollection) (e error) {
	for _,mc := range c.Sensors {
		if mc.AccountId != "" && mc.AssemblyName != "" {
			mkBalance(mc)
			// e = carton.ProvisionerMap[on.Prefix()].TriggerBills(mc.AccountId,mc.AssemblyId, mc.AssemblyName)
			//  if e != nil {
			// 	 return
			//  }
		}
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
		res, e = carton.ProvisionerMap[on.Prefix()].MetricEnvs(time.Now().Add(-10*time.Minute).Unix(),time.Now().Unix(),on.Url,ioutil.Discard)
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
		sc := NewSensor("compute.instance.exists")
		sc.AccountId = h.AccountsId()
		sc.System = on.Prefix()
		sc.Node = h.HostName
		sc.AssemblyId = h.AssemblyId()
		sc.AssemblyName = h.AssemblyName()
		sc.AssembliesId = h.AssembliesId()
		sc.Source = on.Prefix()
		sc.Message = "vm billing"
		sc.Status = h.State()
		sc.AuditPeriodBeginning = time.Unix(metrics.TimeAsInt64(h.VM.Stime), 0).String()
		sc.AuditPeriodEnding = time.Unix(metrics.TimeAsInt64(h.VM.Etime), 0).String()
		sc.AuditPeriodDelta = h.Elapsed()
		sc.addMetric("cpu_cost", h.CpuCost(), h.Cpu(), "delta")
		sc.addMetric("memory_cost", h.MemoryCost(), h.Memory(), "delta")
		sc.addMetric("disk_cost", h.DiskCost(), h.DiskSize(), "delta")
		mc.Add(sc)
	}
	return
}
