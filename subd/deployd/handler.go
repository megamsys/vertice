package deployd

import (
	"github.com/megamsys/megamd/carton"
)

type Handler struct {
	d            *Config
	EventChannel chan bool
}

func NewHandler(c *Config) *Handler {
	return &Handler{d: c}

}

func (h *Handler) serveAMQP(r *carton.Requests) error {
	if _, err := carton.ParseRequest(r.Name); err != nil {
		if rp := carton.NewReqOperator("nil"); rp != nil {
			return rp.Accept(nil)
		}
	}
	return nil
}
