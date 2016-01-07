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

func (h *Handler) serveNSQ(r *carton.Requests) error {
	p, err := carton.ParseRequest(r.CatId, r.Category, r.Action)
	if err != nil {
		return err
	}

	if rp := carton.NewReqOperator(r.CatId); rp != nil {
		return rp.Accept(&p) //error is swalled here.
	}

	return nil
}
