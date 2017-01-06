package metrix

import (
	"github.com/megamsys/vertice/carton"
  "strconv"
	"time"
)

const (
  SNAPSHOTS = "snapshots"
  SNAPSHOT_SENSOR = "instance.snapshots.exists"
)

type Snapshots struct {
	DefaultUnits map[string]string
  RawStatus    []byte
}

type UserSnaps struct {
	AssemblyId    string
	AccountId     string
	AssemblyName  string
	SnapCosts     map[string]string
	TotalStorage  string
}


func (r *Snapshots) Prefix() string {
	return SNAPSHOTS
}

func (r *Snapshots) DeductBill(c *MetricsCollection) (e error) {
	for _, mc := range c.Sensors {
		mkBalance(mc, r.DefaultUnits)
	}
	return
}

func (s *Snapshots) Collect(c *MetricsCollection) (e error) {
  snp := carton.Snaps{}
  snps, e :=  snp.GetBox()
  if e != nil {
		return
	}

	s.CollectMetricsFromStats(c, snps)
	e = s.DeductBill(c)
	return
}

func (c *Snapshots) ReadUsers() ([]carton.Account, error) {
	act := new(carton.Account)
	res, e := act.GetUsers()
	if e != nil {
		return nil, e
	}
	return res, nil
}

//actually the NewSensor can create trypes based on the event type.
func (c *Snapshots) CollectMetricsFromStats(mc *MetricsCollection, snps []carton.Snaps) {
	for _, a := range snps {
		sc := NewSensor(SNAPSHOT_SENSOR)
		sc.AccountId = a.AccountId
		sc.AssemblyId = a.AssemblyId
		sc.System = c.Prefix()
		sc.Node = ""
		sc.AssemblyName = a.Id
		sc.Source = c.Prefix()
		sc.Message = "snapshot billing"
		sc.Status = "health-ok"
		sc.AuditPeriodBeginning = strconv.FormatInt(time.Now().Add(-MetricsInterval).Unix(), 10)
		sc.AuditPeriodEnding = strconv.FormatInt(time.Now().Unix(), 10)
		sc.AuditPeriodDelta = ""
		sc.addMetric(STORAGE_COST, c.DefaultUnits[STORAGE_COST_PER_HOUR], a.Sizeof(), "delta")
		sc.CreatedAt = time.Now()
		mc.Add(sc)
	}

	return
}
