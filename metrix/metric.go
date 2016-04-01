package metrix

import (
	"gopkg.in/yaml.v2"
	"strconv"
)

type Metric struct {
	MetricName  string `json:"metric_name"`
	MetricValue string `json:"metric_value"`
	MetricUnits string `json:"metric_units"`
	MetricType  string `json:"metric_type"`
}

func (m *Metric) String() string {
	if d, err := yaml.Marshal(m); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}
func parseInt(s string) (i int) {
	i, _ = strconv.Atoi(s)
	return
}

func parseInt64(s string) (i int64) {
	i, _ = strconv.ParseInt(s, 10, 64)
	return
}
