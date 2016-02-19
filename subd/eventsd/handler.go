package eventsd

import (
	"github.com/megamsys/vertice/events"
)

type Handler struct {
	d            *Config
	EventChannel chan bool
}

func NewHandler(c *Config) *Handler {
	return &Handler{d: c}

}

func (h *Handler) serveNSQ(e *events.Event) error {
	if err := events.W.Write(e); err != nil {
		return err
	}
	return nil
}
