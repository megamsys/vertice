package docker

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"time"

	"github.com/megamsys/services/deployd"
)

type Handler struct {
	Provider string
	One      *OneConfig
}

func (h *Handler) ServeAMQP(r app.Requests) error {
	assembly, err := app.Get(r.Id)
	if err != nil {
		return err
	}

	di := app.DeployInfo{
		Assembly: assembly,
		Config:   h.One,
		//	EventChannel:      ch,
	}

	switch ParseRequest(r.Name) {
	case ReqCreate:
		if err = provMap[h.Provider].Deploy(di); err != nil {
			return err
		}
	case ReqCreating:
		if err = provMap[h.Provider].Stateup(di); err != nil {
			return err
		}
	case ReqDelete:
		if err = provMap[h.Provider].UnDeploy(di); err != nil {
			return err
		}
	default:
		ReqError
	}
	return nil
}
