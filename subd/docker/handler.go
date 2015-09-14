package docker

import (
	"github.com/megamsys/megamd/carton"
)

type Handler struct {
	Provider string
	D        *Config
}

// NewHandler returns a new instance of handler with routes.
func NewHandler(c *Config) *Handler {
	return &Handler{D: c}
}

func (h *Handler) serveAMQP(r *carton.Requests) error {
	return nil
}
