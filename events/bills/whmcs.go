package bills

import (
	log "github.com/Sirupsen/logrus"
)

const (
	billerName = "whmcs"
)

func init() {
	Register(billerName, createBiller())
}

type whmcsBiller struct {
	enabled bool
	apiKey  string
	domain  string
}

func createBiller() BillProvider {
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

func (w whmcsBiller) Onboard(o *BillOpts) error {
	return nil
}

func (w whmcsBiller) Deduct(o *BillOpts) error {
	return nil
}

func (w whmcsBiller) Transaction(o *BillOpts) error {
	return nil
}

func (w whmcsBiller) Invoice(o *BillOpts) error {
	return nil
}

func (w whmcsBiller) Nuke(o *BillOpts) error {
	return nil
}

func (w whmcsBiller) Suspend(o *BillOpts) error {
	return nil
}

func (w whmcsBiller) Notify(o *BillOpts) error {
	return nil
}
