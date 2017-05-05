package metricsd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/vertice/metrix"
)

type Handler struct {
	EventChannel chan bool
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) processCollector(mh *metrix.MetricHandler, output *metrix.OutputHandler, c metrix.MetricCollector) error {
	all, err := mh.Collect(c)
	if err != nil {
		log.Debugf("%v", err)
		return err
	}
	return output.WriteMetrics(all)
}
