package metrix

import (
	"encoding/json"
	"time"

	log "github.com/Sirupsen/logrus"
)

func SendMetricsToRiak(address []string, metrics []*Metric, hostname string) (e error) {
	_ = time.Now()
	for _, m := range metrics {
		_ = hostname + "." + m.Key
		_, e := json.Marshal(m)
		if e != nil {
			log.Debugf(e.Error())
			continue
		}

		/*if e != nil {
			log.Debugf(e.Error())
			continue
		}*/
	}
	//fmt.Printf("sent %d metrics in %.06f\n", len(metrics), time.Now().Sub(started).Seconds())
	return nil
}
