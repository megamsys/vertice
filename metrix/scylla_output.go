package metrix

import (
	log "github.com/Sirupsen/logrus"
	ldb "github.com/megamsys/libgo/db"
	"github.com/megamsys/libgo/events"
	"github.com/megamsys/libgo/events/alerts"
	constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/vertice/meta"
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
		ops := ldb.Options{
			TableName:   SENSORSBUCKET,
			Pks:         []string{"sensor_type"},
			Ccms:        []string{"account_id", "assembly_id"},
			Hosts:       meta.MC.Scylla,
			Keyspace:    meta.MC.ScyllaKeyspace,
			Username:    meta.MC.ScyllaUsername,
			Password:    meta.MC.ScyllaPassword,
			PksClauses:  map[string]interface{}{"sensor_type": m.SensorType},
			CcmsClauses: map[string]interface{}{"account_id": m.AccountId, "assembly_id": m.AssemblyId},
		}
		s := m.ParseScyllaformat()
		if err = ldb.Storedb(ops, s); err != nil {
			log.Debugf(err.Error())
			continue
		}
		//make it posted in a background lowpriority channel
		//mkBalance(m)
	}
	log.Debugf("sent %d metrics in %.06f\n", len(metrics), time.Since(started).Seconds())
	return nil
}

func mkBalance(s *Sensor, du map[string]string) error {
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
