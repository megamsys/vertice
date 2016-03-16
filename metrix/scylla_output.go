package metrix

import (
	"time"

	log "github.com/Sirupsen/logrus"
	ldb "github.com/megamsys/libgo/db"
	"github.com/megamsys/libgo/events"
	"github.com/megamsys/libgo/events/alerts"
	"github.com/megamsys/vertice/meta"
)

func SendMetricsToScylla(address []string, metrics Sensors, hostname string) (err error) {
	started := time.Now()
	for _, m := range metrics {
		ops := ldb.Options{
			TableName:   SENSORSBUCKET,
			Pks:         []string{"type"},
			Ccms:        []string{"account_id"},
			Hosts:       meta.MC.Scylla,
			Keyspace:    meta.MC.ScyllaKeyspace,
			PksClauses:  map[string]interface{}{"type": m.SensorType},
			CcmsClauses: map[string]interface{}{"account_id": m.AccountsId},
		}
		s, err := m.ParseScyllaformat()
		if err != nil {
			log.Debugf(err.Error())
			continue
		}
		if err = ldb.Storedb(ops, s); err != nil {
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
				EventData:   alerts.EventData{M: mi},
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
				EventData:   alerts.EventData{M: mi},
				Timestamp:   time.Now().Local(),
			},
		})
	return newEvent.Write()
}
