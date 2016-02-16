package metrix

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/vertice/db"
	"github.com/megamsys/vertice/events"
	"github.com/megamsys/vertice/events/alerts"
)

func SendMetricsToScylla(address []string, metrics Sensors, hostname string) (err error) {
	started := time.Now()
	for _, m := range metrics {
		if err = db.Store(SENSORSBUCKET, m.Id, m); err != nil {
			log.Debugf(err.Error())
			continue
		}
		//make it posted in a background lowpriority channel
		mkBalance(m)
		mkTransaction(m)
	}
	log.Debugf("sent %d metrics in %.06f\n", len(metrics), time.Since(started).Seconds())
	return nil
}

func mkBalance(s *Sensor) error {
	mi := make(map[string]string)
	mi[alerts.VERTNAME] = ""
	mi[alerts.COST] = "0.1" //stick the cost from metrics
	newEvent := events.NewMulti(
		[]*events.Event{
			&events.Event{
				AccountsId:  s.AccountsId,
				EventAction: alerts.DEDUCT,
				EventType:   events.EventBill,
				EventData:   events.EventData{M: mi},
				Timestamp:   time.Now().Local(),
			},
		})
	return newEvent.Write()
}

func mkTransaction(s *Sensor) error {
	mi := make(map[string]string)
	mi[alerts.VERTNAME] = ""
	mi[alerts.COST] = "0.1" //stick the cost from metrics
	newEvent := events.NewMulti(
		[]*events.Event{
			&events.Event{
				AccountsId:  s.AccountsId,
				EventAction: alerts.TRANSACTION, //Change type to transaction
				EventType:   events.EventBill,
				EventData:   events.EventData{M: mi},
				Timestamp:   time.Now().Local(),
			},
		})
	return newEvent.Write()
}
