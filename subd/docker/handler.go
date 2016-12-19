package docker

import (
	"github.com/megamsys/vertice/carton"
)

type Handler struct {
	Provider string
	D        *Config
}

// NewHandler returns a new instance of handler with routes.
func NewHandler(c *Config) *Handler {
	return &Handler{D: c}
}

func (h *Handler) serveNSQ(r *carton.Requests) error {
	p, err := carton.ParseRequest(r)
	if err != nil {
		return err
	}

	if rp := carton.NewReqOperator(r); rp != nil {
		return rp.Accept(&p) //error is swalled here.
	}
	return nil
}
