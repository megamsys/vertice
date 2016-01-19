package metrix

import (
	"encoding/json"
	"time"

	log "github.com/Sirupsen/logrus"
)

func SendMetricsToRiak(address []string, metrics []*Metric, hostname string) (err error) {
	started := time.Now()
	for _, m := range metrics {
		_, err := json.Marshal(m)
		if err != nil {
			log.Debugf(err.Error())
			continue
		}

		if err != nil {
			log.Debugf(err.Error())
			continue
		}
	}
	log.Debugf("sent %d metrics in %.06f\n", len(metrics), time.Since(started).Seconds())
	return nil
}
