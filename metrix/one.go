package metrix

import (
	"encoding/xml"
	constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/opennebula-go/metrics"
	"github.com/megamsys/vertice/carton"
	"io/ioutil"
	"strconv"
	"time"
)

const (
	OPENNEBULA = "one"
	QUOTA      = "quota"
	ONDEMAND   = "ondemand"
)

type OpenNebula struct {
	Url          string
	Region       string
	DefaultUnits map[string]string
	SkewsActions map[string]string
	RawStatus    []byte
}

func (on *OpenNebula) Prefix() string {
	return OPENNEBULA
}

func (on *OpenNebula) DeductBill(c *MetricsCollection) (e error) {
	for _, mc := range c.Sensors {
		on.SkewsActions[constants.RESOURCES] = mc.Resources
		if mc.Resources != "" {
			 mkBalance(mc, on.DefaultUnits)
		}

		if on.SkewsActions[constants.ENABLED] == constants.TRUE {
			e = eventSkews(mc, on.SkewsActions)
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
		res, e = carton.ProvisionerMap[on.Prefix()].MetricEnvs(time.Now().Add(-MetricsInterval).Unix(), time.Now().Unix(), on.Region, ioutil.Discard)
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
		var cpu, ram, storage string
		disks := h.Disks()
		cpu = h.VCpu()
		ram = h.Memory()
		sc := NewSensor(ONE_VM_SENSOR)
		if len(h.QuotaId()) > 0 {
			cpu = "0"
			ram = "0"
			if len(disks) > 1 {
				sc.Resources = "disks"
				storage = strconv.FormatInt(on.quotaDisks(disks), 10)
			}
			sc.QuotaId = h.QuotaId()
		} else {
			sc.Resources = "cpu.ram.disks"
			storage = strconv.FormatInt(on.vmDiskUsage(disks), 10)
		}

		sc.AccountId = h.AccountsId()
		sc.System = on.Prefix()
		sc.Node = h.HostName
		sc.AssemblyId = h.AssemblyId()
		sc.AssemblyName = h.AssemblyName()
		sc.AssembliesId = h.AssembliesId()
		sc.Source = on.Prefix()
		sc.Message = "vm billing for " + sc.Resources
		sc.Status = h.State()
		sc.AuditPeriodBeginning = time.Now().Add(-MetricsInterval).Format(time.RFC3339) // time.Unix(h.PStime, 0).String()
		sc.AuditPeriodEnding = time.Now().Format(time.RFC3339)                          // time.Unix(h.PEtime, 0).String()
		sc.AuditPeriodDelta = h.Elapsed()
		sc.addMetric(CPU_COST, h.CpuCost(), cpu, "delta")
		sc.addMetric(MEMORY_COST, h.MemoryCost(), ram, "delta")
		sc.addMetric(DISK_COST, h.DiskCost(), storage, "delta")
		sc.CreatedAt = time.Now()
		if sc.isOk() {
			mc.Add(sc)
		}
	}
	return
}

func (on *OpenNebula) vmDiskUsage(disks []metrics.Disk) int64 {
	var totalsize int64
	for _, v := range disks {
		totalsize = totalsize + v.Size
	}
	return totalsize
}

func (on *OpenNebula) quotaDisks(disks []metrics.Disk) int64 {
	return on.vmDiskUsage(disks) - disks[0].Size
}
