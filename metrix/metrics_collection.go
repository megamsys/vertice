package metrix

type MetricsCollection struct {
	Prefix  string
	Metrics Metrics
	Mapping map[string]string
}

func (m *MetricsCollection) AddSingularMappings(mappings []string) {
	for _, k := range mappings {
		m.AddSingularMapping(k)
	}
}

func (m *MetricsCollection) AddSingularMapping(from string) {
	m.AddMapping(from, from)
}

func (m *MetricsCollection) AddMapping(from, to string) {
	if m.Mapping == nil {
		m.Mapping = map[string]string{}
	}
	m.Mapping[from] = to
}

func (m *MetricsCollection) AddWithTags(key string, v string, tags map[string]string) (e error) {
	if realKey, ok := m.Mapping[key]; ok {
		key = realKey
	}
	m.Metrics = append(m.Metrics, &Metric{Key: m.Prefix + "." + key, Value: v, Tags: tags})
	return
}

func (m *MetricsCollection) MustAddString(key, value string) error {
	return m.Add(key, value)
}

func (m *MetricsCollection) Add(key string, v string) (e error) {
	return m.AddWithTags(key, v, map[string]string{})
}
