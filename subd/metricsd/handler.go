package metricsd

import (
	"github.com/megamsys/vertice/metrix"
)

type Handler struct {
	EventChannel chan bool
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) processCollector(mh *metrix.MetricHandler, output *metrix.OutputHandler, c metrix.MetricCollector) error {
	if _, err := mh.Collect(c); err != nil {
		return err
	}
	return output.WriteMetrics(all)
}
