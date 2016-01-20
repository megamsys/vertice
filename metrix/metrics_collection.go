package metrix

import "gopkg.in/yaml.v2"

type Sensors []*Sensor

type MetricsCollection struct {
	Prefix  string
	Sensors Sensors
}

func (m *MetricsCollection) Add(s *Sensor) {
	m.Sensors = append(m.Sensors, s)
	return
}

func (m *MetricsCollection) String() string {
	if d, err := yaml.Marshal(m); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}
