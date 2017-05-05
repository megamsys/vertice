package metrix

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/api"
	"github.com/megamsys/libgo/events"
	"github.com/megamsys/libgo/events/alerts"
	constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/vertice/carton"
	"strconv"
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
	cb, _ := strconv.ParseFloat(m, 64)
	if cb <= 0 {
		return nil
	}
	mi[constants.ACCOUNTID] = s.AccountId
	mi[constants.ASSEMBLYID] = s.AssemblyId
	mi[constants.ASSEMBLIESID] = s.AssembliesId
	mi[constants.ASSEMBLYNAME] = s.AssemblyName
	mi[constants.RESOURCES] = s.Resources
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

func eventSkews(s *Sensor, action alerts.EventAction, skews map[string]string) error {
	mi := make(map[string]string, 0)
	mi[constants.ACCOUNTID] = s.AccountId
	mi[constants.ASSEMBLYID] = s.AssemblyId
	mi[constants.ASSEMBLIESID] = s.AssembliesId
	mi[constants.ASSEMBLYNAME] = s.AssemblyName
	mi[constants.QUOTAID] = s.QuotaId

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
