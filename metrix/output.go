package metrix

import (
	"fmt"
	"os"
	"time"
)

type OutputHandler struct {
	RiakAddress []string
	Hostname    string
}

func (o *OutputHandler) WriteMetrics(all []*Metric) (e error) {
	if o.Hostname == "" {
		hn, e := os.Hostname()
		if e != nil {
			return e
		}
		o.Hostname = hn
	}
	sent := false
	if len(o.RiakAddress) > 0 {
		e = SendMetricsToRiak(o.RiakAddress, all, o.Hostname)
		sent = true
	}

	if !sent {
		SendMetricsToStdout(all, o.Hostname)
	}
	return
}

func SendMetricsToStdout(metrics []*Metric, hostname string) {
	now := time.Now()
	for _, m := range metrics {
		fmt.Println(m.Ascii(now, hostname))
	}
}
