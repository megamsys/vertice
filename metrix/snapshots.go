package metrix

import (
  "time"
	"github.com/megamsys/vertice/snapshots"
	"github.com/megamsys/vertice/carton"
)

const CEPHRBD = "ceph_rbd"

type Snapshot struct {
	CephRbd CephRbd
}

type CephRbd struct {
	Server   		 Server
	PoolName 		 string
	DefaultUnits map[string]string
}

type Server struct {
	Host             string            `json:"ipaddress"`
	Username         string            `json:"username"`
	Password         string            `json:"password"`
	PrivateKey       string            `json:"privatekey"`
}

func (r *CephRbd) Prefix() string {
	return CEPHRBD
}

func (r *CephRbd) DeductBill(c *MetricsCollection) (e error) {
	for _, mc := range c.Sensors {
		mkBalance(mc, r.Inputs)
	}
	return
}

func (r *CephRbd) Collect(c *MetricsCollection) (e error) {
	acc, e := r.ReadUsers()
	if e != nil {
		return
	}

  rbd := &snapshots.CephRbdRunner{
		Host: r.Server.Host,
		Username: r.Server.Username,
		Password: r.Server.Password,
		PrivateKey: r.Server.PrivateKey,
	}

	snps := rbd.GetUserSnaps(acc)
	if len(snps) < 1 {
		return
	}

	r.CollectMetricsFromStats(c, snps)
	e = r.DeductBill(c)
	return
}

func (c *CephRbd) ReadUsers() ([]carton.Account, error) {
	act := new(carton.Account)
	res, e := act.GetUsers()
	if e != nil {
		return nil, e
	}
	return res, nil
}

//actually the NewSensor can create trypes based on the event type.
func (c *CephRbd) CollectMetricsFromStats(mc *MetricsCollection, snps []snapshots.AsmSnaps) {
	for _, a := range snp {
			sc := NewSensor(CEPH_SNAPSHOT_SENSOR)
			sc.AccountId = a.AccountId
			sc.AssemblyId = a.AssemblyId
			sc.System = c.Prefix()
			sc.Node = c.Host
			sc.AssemblyName = a.AssemblyName
			sc.Source = c.Prefix()
			sc.Message = "snapshot billing"
			sc.Status = "health-ok"
			sc.AuditPeriodBeginning = time.Now().Add(-10 * time.Minute).String()
			sc.AuditPeriodEnding = time.Now().String()
			sc.AuditPeriodDelta = ""
			sc.addMetric(STORAGE_COST, c.DefaultUnits[STORAGE_COST_PER_HOUR], strconv.FormatFloat(r.TotalSizeMB, 'f', 4, 64), "delta")
			sc.CreatedAt = time.Now()
			mc.Add(sc)
	}

	return
}
