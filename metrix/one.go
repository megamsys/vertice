package metrix

import (
	"encoding/json"
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
		on.RawStatus, e = FetchURL(on.Url)
		if e != nil {
			return
		}
	}
	b = on.RawStatus
	return
}

func (on *OpenNebula) ParseStatus(b []byte) (ons *OpenNebulaStatus, e error) {
	ons = &OpenNebulaStatus{}
	e = json.Unmarshal(b, ons)
	if e != nil {
		return nil, e
	}
	return ons, nil
}

func (on *OpenNebula) CollectMetricsFromStats(mc *MetricsCollection, s *OpenNebulaStatus) {
	for id, h := range s.HISTORYS {
		tags := map[string]string{"machine_id": string(id)}
		mc.AddWithTags("node", h.HOSTNAME, tags)
		mc.AddWithTags("accounts_id", h.AccountsId(), tags)
		mc.AddWithTags("assembly_id", h.AssemblyId(), tags)
		mc.AddWithTags("assembly_name", h.AssemblyName(), tags)
		mc.AddWithTags("assemblies_id", h.AssembliesId(), tags)
		mc.AddWithTags("status", h.State(), tags)
		mc.AddWithTags("system", "one", tags)
		mc.AddWithTags("cpu", h.Cpu(), tags)
		mc.AddWithTags("memory", h.Memory(), tags)
		mc.AddWithTags("cpu_cost", h.CpuCost(), tags)
		mc.AddWithTags("memory_cost", h.MemoryCost(), tags)
		mc.AddWithTags("audit_period_beginning", time.Unix(timeAsInt64(h.VM.STIME), 0).String(), tags) // UTC
		mc.AddWithTags("audit_period_ending", time.Unix(timeAsInt64(h.VM.ETIME), 0).String(), tags)    //UTC
		mc.AddWithTags("audit_period_delta", h.VM.elapsed(), tags)                                     //Hours
	}
	return
}
