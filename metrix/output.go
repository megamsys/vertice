package metrix

import (
	"fmt"
	"os"
)

type OutputHandler struct {
	ScyllaAddress string
	Hostname      string
}

func (o *OutputHandler) WriteMetrics(all Sensors) (e error) {
	if o.Hostname == "" {
		hn, e := os.Hostname()
		if e != nil {
			return e
		}
		o.Hostname = hn
	}
	sent := false
	if len(o.ScyllaAddress) > 0 {
		e = SendMetricsToScylla(all, o.Hostname)
		sent = true
	}

	if !sent {
		SendMetricsToStdout(all, o.Hostname)
	}
	return
}

func SendMetricsToStdout(metrics Sensors, hostname string) {
	for _, m := range metrics {
		fmt.Println(m)
	}
}
