package marketplacesd

import (
 	"github.com/megamsys/vertice/marketplaces"
 	log "github.com/Sirupsen/logrus"
	"fmt"
 )

type Handler struct {
	d            *Config
	EventChannel chan bool
}

func NewHandler(c *Config) *Handler {
	return &Handler{d: c}

}

func (h *Handler) serveNSQ(r *marketplaces.Requests) error {
	req, err := r.ParseRequest()
	if err != nil {
		log.Errorf("Error parsing request : %s  -  %s  : %s",r.Category, r.Action, err)
		return err
	}
	fmt.Println("*************************",req)
  return req.Process(r.Action)
}
