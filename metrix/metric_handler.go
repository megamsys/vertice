package metrix

type MetricHandler struct {
}

type MetricCollector interface {
	Prefix() string
	Collect(*MetricsCollection) error
}

func (h *MetricHandler) Collect(c MetricCollector) (Sensors, error) {
	sc := &MetricsCollection{Prefix: c.Prefix()}
	e := c.Collect(sc)
	return sc.Sensors, e
}
