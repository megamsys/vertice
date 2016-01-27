package eventsd

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

func (h *Handler) serveNSQ(r *carton.Requests) error {
	return nil
}
