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

func mkBalance(s *Sensor, du map[string]string) error {
	mi := make(map[string]string, 0)
	m := s.Metrics.Totalcost(du)
	mi[constants.ACCOUNTID] = s.AccountId
	mi[constants.ASSEMBLYID] = s.AssemblyId
	mi[constants.ASSEMBLYNAME] = s.AssemblyName
	mi[constants.CONSUMED] = m
	mi[constants.START_TIME] = s.AuditPeriodBeginning
	mi[constants.END_TIME] = s.AuditPeriodEnding
	mi[constants.BILL_TYPE] = s.SensorType

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

func eventSkews(s *Sensor, skews map[string]string) error {
	var action alerts.EventAction
	mi := make(map[string]string, 0)

	mi[constants.ACCOUNTID] = s.AccountId
	mi[constants.ASSEMBLYID] = s.AssemblyId
	mi[constants.ASSEMBLYNAME] = s.AssemblyName
	mi[constants.QUOTAID] = s.QuotaId
	if len(s.QuotaId) > 0 {
		quota, err := carton.NewQuota(s.AccountId, s.QuotaId)
		if err != nil {
			return err
		}
		action = alerts.QUOTA_UNPAID
		if quota.Status == "paid" {
			return nil
		}
		mi[constants.SKEWS_TYPE] = "vm.quota.unpaid"
	} else {
		action = alerts.SKEWS_ACTIONS
		mi[constants.SKEWS_TYPE] = "vm.ondemand.bills"
	}

	for k, v := range skews {
		mi[k] = v
	}

	newEvent := events.NewMulti(
		[]*events.Event{
			&events.Event{
				AccountsId:  s.AccountId,
				EventAction: action,
				EventType:   constants.EventBill,
				EventData:   alerts.EventData{M: mi},
				Timestamp:   time.Now().Local(),
			},
		})
	return newEvent.Write()
}
