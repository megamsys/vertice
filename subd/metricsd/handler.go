package metricsd

import (
	"github.com/megamsys/megamd/metrix"
)

type Handler struct {
	EventChannel chan bool
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) processCollector(mh *metrix.MetricHandler,
	output *metrix.OutputHandler, c metrix.MetricCollector) error {

	all, err := mh.Collect(c)
	if err != nil {
		return err
	}
	if err = output.WriteMetrics(all); err != nil {
		return err
	}
	return nil
}
