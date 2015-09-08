package deployd

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/megamd/carton"
)

type Handler struct {
	d            *Config
	EventChannel chan bool
}

// NewHandler returns a new instance of handler with routes.
func NewHandler(c *Config) *Handler {
	return &Handler{d: c}

}

func (h *Handler) serveAMQP(r *carton.Requests) error {
	a, err := carton.Get(r.CatId)
	if err != nil {
		return err
	}

	c, err := a.MkCartons()
	if err != nil {
		return err
	}

	ra, err := carton.ParseRequest(r.Name)
	log.Warnf("Parse request %s\n", r.Name)
	switch ra {
	case carton.ReqCreate:
		carton.CDeploys(c[0])
		fmt.Printf("%s", "OK")
	case carton.ReqDelete:
		carton.Delete(c[0])
		/*case ReqCreating:
		if err = app.Stateup(di); err != nil {
		}
			return err
		*/
	default:
		return carton.ErrInvalidReqtype
	}
	return nil
}
