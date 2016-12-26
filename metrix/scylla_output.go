package metrix

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/api"
	"github.com/megamsys/libgo/events"
	"github.com/megamsys/libgo/events/alerts"
	constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/vertice/carton"
	"time"
)

func SendMetricsToScylla(metrics Sensors, hostname string) (err error) {
	started := time.Now()
	for _, m := range metrics {
			cl := api.NewClient(carton.NewArgs(m.AccountId, ""), "/sensors/content")
			if _, err := cl.Post(m); err != nil {
				log.Debugf(err.Error())
				continue
			}
	}
	log.Debugf("sent %d metrics in %.06f\n", len(metrics), time.Since(started).Seconds())
	return nil
}

func quotaChecker(id, email string) (bool, error) {
	asm, err := carton.NewAssembly(id, email, "")
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

	if flag, err := quotaChecker(s.AssemblyId, s.AccountId); !flag {
		return err
	}
	mi := make(map[string]string)
	m := s.Metrics.Totalcost(du)
	mi[constants.ACCOUNTID] = s.AccountId
	mi[constants.ASSEMBLYID] = s.AssemblyId
	mi[constants.ASSEMBLYNAME] = s.AssemblyName
	mi[constants.CONSUMED] = m
	mi[constants.START_TIME] = s.AuditPeriodBeginning
	mi[constants.END_TIME] = s.AuditPeriodEnding

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
				EventAction: alerts.BILLEDHISTORY, //Change type to transaction
				EventType:   constants.EventBill,
				EventData:   alerts.EventData{M: mi},
				Timestamp:   time.Now().Local(),
			},
		})
	return newEvent.Write()
}
