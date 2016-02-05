package metrix

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/vertice/db"
)

func SendMetricsToRiak(address []string, metrics Sensors, hostname string) (err error) {
	started := time.Now()
	for _, m := range metrics {
		if err = db.Store(SENSORSBUCKET, m.Id, m); err != nil {
			log.Debugf(err.Error())
			continue
		}
	}
	log.Debugf("sent %d metrics in %.06f\n", len(metrics), time.Since(started).Seconds())
	return nil
}
