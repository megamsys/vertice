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

func (h *Handler) processCollector(mh *metrix.MetricHandler,
	output *metrix.OutputHandler, c metrix.MetricCollector) error {

	_, err := mh.Collect(c)
	if err != nil {
		return err
	}
	// if err = output.WriteMetrics(all); err != nil {
	// 	return err
	// }
	return nil
}
