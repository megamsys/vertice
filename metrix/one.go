package metrix

import (
	"encoding/xml"
	"io/ioutil"
	"time"

	"github.com/megamsys/vertice/carton"
)

const OPENNEBULA = "one"

type OpenNebula struct {
	Url       string
	RawStatus []byte
}

func (on *OpenNebula) Prefix() string {
	return "one"
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
	return
}

func (on *OpenNebula) ReadStatus() (b []byte, e error) {
	if len(on.RawStatus) == 0 {
		var res []interface{}
		res, e = carton.ProvisionerMap[on.Prefix()].MetricEnvs(time.Now().Add(-10*time.Minute).Unix(),
			time.Now().Unix(), ioutil.Discard)
		if e != nil {
			return
		}
		on.RawStatus = []byte(res[1].(string))
	}
	b = on.RawStatus
	return
}

func (on *OpenNebula) ParseStatus(b []byte) (ons *OpenNebulaStatus, e error) {
	ons = &OpenNebulaStatus{}
	e = xml.Unmarshal(b, ons)
	if e != nil {
		return nil, e
	}
	return ons, nil
}

//actually the NewSensor can create trypes based on the event type.
func (on *OpenNebula) CollectMetricsFromStats(mc *MetricsCollection, s *OpenNebulaStatus) {
	for _, h := range s.HISTORYS {
		sc := NewSensor("compute.instance.exists")
		sc.AccountsId = h.AccountsId()
		sc.addPayload(&Payload{System: on.Prefix(),
			Node:         h.HOSTNAME,
			AssemblyId:   h.AssemblyId(),
			AssemblyName: h.AssemblyName(),
			AssembliesId: h.AssembliesId(),
			Source:       on.Prefix(),
			Message:      "vm billing",
			Status:       h.State(),
			BeginAudit:   time.Unix(timeAsInt64(h.VM.STIME), 0).String(),
			EndAudit:     time.Unix(timeAsInt64(h.VM.ETIME), 0).String(),
			DeltaAudit:   h.VM.elapsed()})

		sc.addMetric("cpu_cost", h.CpuCost(), h.Cpu(), "delta")
		sc.addMetric("memory_cost", h.MemoryCost(), h.Memory(), "delta")
		mc.Add(sc)
	}
	return
}
