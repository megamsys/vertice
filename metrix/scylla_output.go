package metrix

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/events"
	"github.com/megamsys/vertice/carton"
	"github.com/megamsys/libgo/events/alerts"
	constants "github.com/megamsys/libgo/utils"
	 "github.com/megamsys/libgo/api"
	"time"
)

const (
	ACCOUNTID    = "AccountId"
	ASSEMBLYID   = "AssemblyId"
	ASSEMBLYNAME = "AssemblyName"
	CONSUMED     = "Consumed"
	STARTTIME    = "StartTime"
	ENDTIME      = "EndTime"
)

func SendMetricsToScylla(address []string, metrics Sensors, hostname string) (err error) {
	started := time.Now()
	for _, m := range metrics {
		args := carton.NewArgs(m.AccountId, "")
		s := m.ParseScyllaformat()
		args.Path = "/sensors/content"
		cl := api.NewClient(args)
		if _, err := cl.Post(s); err != nil {
			log.Debugf(err.Error())
			continue
		}
		//make it posted in a background lowpriority channel
		//mkBalance(m)
	}
	log.Debugf("sent %d metrics in %.06f\n", len(metrics), time.Since(started).Seconds())
	return nil
}

func quotaChecker(id, email string) (bool, error) {
	asm, err := carton.NewAssembly(id,email,"")
	if err != nil {
		return false, err
	}

  qid := asm.QuotaID()
	if len(qid) > 0 {
		return false, nil
	}

	return true, nil

}

func mkBalance(s *Sensor, du map[string]string) error {

	if flag, err := quotaChecker(s.AssemblyId,s.AccountId); !flag {
		return err
	}

	mi := make(map[string]string)
	m := s.Metrics.Totalcost(du)
	mi[ACCOUNTID] = s.AccountId
	mi[ASSEMBLYID] = s.AssemblyId
	mi[ASSEMBLYNAME] = s.AssemblyName
	mi[CONSUMED] = m //stick the cost from metrics
	mi[STARTTIME] = s.AuditPeriodBeginning
	mi[ENDTIME] = s.AuditPeriodEnding

	newEvent := events.NewMulti(
		[]*events.Event{
			&events.Event{
				AccountsId:  s.AccountId,
				EventAction: alerts.DEDUCT,
				EventType:   constants.EventBill,
				EventData:   alerts.EventData{M: mi},
				Timestamp:   time.Now().Local(),
			},
			&events.Event{
				AccountsId:  s.AccountId,
				EventAction: alerts.TRANSACTION, //Change type to transaction
				EventType:   constants.EventBill,
				EventData:   alerts.EventData{M: mi},
				Timestamp:   time.Now().Local(),
			},
		})
	return newEvent.Write()
}
