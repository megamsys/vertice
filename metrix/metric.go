package metrix

import (
	"strconv"

	"gopkg.in/yaml.v2"
)

type Metric struct {
	Key   string `json:"name"`
	Value string `json:"value"`
	Units string `json:"units"`
	Type  string `json:"type"`
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
