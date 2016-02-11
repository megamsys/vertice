package bills

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/vertice/events"
)

const (
	billerName = "whmcs"
)

func init() {
	events.Register(billerName, createBiller())
}

type whmcsBiller struct {
	enabled bool
	apiKey  string
	domain  string
}

func createBiller() events.BillProvider {
	vBiller := whmcsBiller{
		enabled: false,
		apiKey:  "",
		domain:  "",
	}
	log.Debugf("%s ready", billerName)
	return vBiller
}

func (w *whmcsBiller) String() string {
	return "WHMCS:(" + w.apiKey + "," + w.domain + ")"
}

func (w whmcsBiller) IsEnabled() bool {
	return w.enabled
}

func (w whmcsBiller) Onboard(o *events.BillOpts) error {
	return nil
}

func (w whmcsBiller) Deduct(o *events.BillOpts) error {
	return nil
}

func (w whmcsBiller) Invoice(o *events.BillOpts) error {
	return nil
}

func (w whmcsBiller) Nuke(o *events.BillOpts) error {
	return nil
}

func (w whmcsBiller) Suspend(o *events.BillOpts) error {
	return nil
}

func (w whmcsBiller) Notify(o *events.BillOpts) error {
	return nil
}
