package metrix

import (
	"encoding/json"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/megamd/db"
)

func SendMetricsToRiak(address []string, metrics Sensors, hostname string) (err error) {
	started := time.Now()
	for _, m := range metrics {
		sj, err := json.Marshal(m)
		if err != nil {
			log.Debugf(err.Error())
			continue
		}

		if err = db.Store(SENSORSBUCKET, m.Id, sj); err != nil {
			log.Debugf(err.Error())
			continue
		}
	}
	log.Debugf("sent %d metrics in %.06f\n", len(metrics), time.Since(started).Seconds())
	return nil
}
