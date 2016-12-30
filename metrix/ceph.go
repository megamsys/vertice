package metrix
/*
import (
	// "encoding/xml"
	// "github.com/megamsys/vertice/carton"
	// "io/ioutil"
	// "strconv"
	"time"
)

const (
	CEPHRGW  = "ceph_rgw"
  CEPHRBD  = "ceph_rbd"
)

type CephStorage struct {
	Url          string
  AccessKey    string
  SecretKey    string
  UserId       string
  UserPrefix   string
	DefaultUnits map[string]string
	RawStatus    []byte
}

type CephRGWStats struct {

}

func (rgw *CephStorage) Prefix() string {
	return "ceph_rgw"
}

func (rgw *CephStorage) DeductBill(c *MetricsCollection) (e error) {
	for _, mc := range c.Sensors {
			mkBalance(mc, rgw.DefaultUnits)
	}
	return
}

func (rgw *CephStorage) Collect(c *MetricsCollection) (e error) {
	b, e := rgw.ReadStatus()
	if e != nil {
		return
	}
	 e = rgw.ParseStatus(b)
	if e != nil {
		return
	}
	//rgw.CollectMetricsFromStats(c, s)
	e = rgw.DeductBill(c)
	return
}

func (rgw *CephStorage) ReadStatus() (b []byte, e error) {
	if len(rgw.RawStatus) == 0 {
		// var res []interface{}
		rgw.RawStatus = []byte("result")
	}

	b = rgw.RawStatus
	return
}

func (rgw *CephStorage) ParseStatus(b []byte) (e error) {

	return  nil
}

//actually the NewSensor can create trypes based on the event type.
func (rgw *CephStorage) CollectMetricsFromStats(mc *MetricsCollection, s *CephRGWStats) {
  //for _, h := range s.Stats { }
		sc := NewSensor("compute.instance.exists")
		sc.AccountId = ""
		sc.System = rgw.Prefix()
		sc.Node = ""
		sc.AssemblyId = ""
		sc.AssemblyName = ""
		sc.AssembliesId = ""
		sc.Source = rgw.Prefix()
		sc.Message = "vm billing"
		sc.Status = ""
		sc.AuditPeriodBeginning = "time.Unix(h.PStime, 0).String()"
		sc.AuditPeriodEnding = "time.Unix(h.PEtime, 0).String()"
		sc.AuditPeriodDelta = ""
		sc.addMetric(CPU_COST, "CpuCost", "used cup", "delta")
		sc.addMetric(MEMORY_COST, "MemoryCost", "used Memory", "delta")
		sc.addMetric(DISK_COST, "DiskCost","used disk", "delta")
		sc.CreatedAt = time.Now()
		if sc.isBillable() {
			mc.Add(sc)
		}
	return
}
*/
