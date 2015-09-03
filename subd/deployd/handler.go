package deployd

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"time"

)

type Handler struct {
	d  *deployd.Config
	ch EventChannel
}

func (h *Handler) ServeAMQP(r app.Requests) error {
	c, err := carton.Get(r.Id)
	if err != nil {
		return err
	}

	switch ParseRequest(r.Name) {
	case ReqCreate:
		err = carton.Deploy(app.DeployOptions{
			C:      c,
			Image:  c.Image,
			Config: h.d,
		})
		if err == nil {
			fmt.Fprintln(w, "\nOK")
		}
	case ReqDelete:
		if err = carton.Delete(dd); err != nil {
			return err
		}
		/*case ReqCreating:
		if err = app.Stateup(di); err != nil {
		}
			return err
		*/
	default:
		ReqError
	}
	return nil
}
