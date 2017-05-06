package metrix

import (
	"github.com/megamsys/vertice/carton"
	"github.com/megamsys/vertice/storage"
	"strconv"
	"time"
)

const (
	CEPHRGW = "ceph_rgw"
)

type CephStorage struct {
}

type CephRGWStats struct {
	Url          string
	AdminUser    string
	MasterKey    string
	AccessKey    string
	SecretKey    string
	UserId       string
	UserPrefix   string
	DefaultUnits map[string]string
	RawStatus    []byte
}

func (rgw *CephRGWStats) Prefix() string {
	return CEPHRGW
}

func (rgw *CephRGWStats) DeductBill(c *MetricsCollection) (e error) {
	for _, mc := range c.Sensors {
		mkBalance(mc, rgw.DefaultUnits)
	}
	return
}

func (rgw *CephRGWStats) Collect(c *MetricsCollection) (e error) {
	acc, e := rgw.ReadUsers()
	if e != nil {
		return
	}

	rgw.CollectMetricsFromStats(c, acc)
	e = rgw.DeductBill(c)
	return
}

func (c *CephRGWStats) ReadUsers() ([]*carton.Account, error) {
	act := new(carton.Account)
	res, e := act.GetUsers()
	if e != nil {
		return nil, e
	}
	return res, nil
}

//actually the NewSensor can create trypes based on the event type.
func (c *CephRGWStats) CollectMetricsFromStats(mc *MetricsCollection, acts []*carton.Account) {
	for _, a := range acts {
		r := storage.NewRgW(c.Url, c.AccessKey, c.SecretKey)
		r.UserId = a.Email
		err := r.GetUserStorageSize()
		if err == nil {
			sc := NewSensor(CEPH_STORAGE_SENSOR)
			sc.AccountId = a.Email
			sc.System = c.Prefix()
			sc.Node = c.Url
			sc.AssemblyName = ""
			sc.Source = c.Prefix()
			sc.Message = "storage billing"
			sc.Status = "health-ok"
			sc.AuditPeriodBeginning = time.Now().Add(-MetricsInterval).Format(time.RFC3339)
			sc.AuditPeriodEnding = time.Now().Format(time.RFC3339)
			sc.AuditPeriodDelta = ""
			sc.addMetric(STORAGE_COST, c.DefaultUnits[STORAGE_COST_PER_HOUR], strconv.FormatFloat(r.TotalSizeMB, 'f', 4, 64), "delta")
			sc.CreatedAt = time.Now()
			mc.Add(sc)
		}

	}
	return
}
