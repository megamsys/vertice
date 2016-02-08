package bills

import (
	"github.com/megamsys/vertice/events"
)

func init() {
	events.Register("whmcs", whmcsManager{})
}

type whmcsManager struct{}

type WHMCSClient struct{}

func (w whmcsManager) client() (*WHMCSClient, error) {
	//	return whmcs.NewClient(url, tok), nil
	return nil, nil
}

func (w whmcsManager) IsEnabled(o *events.BillOpts) bool {
	return false
}

func (w whmcsManager) Onboard(o *events.BillOpts) error {
	return nil
}

func (w whmcsManager) Deduct(o *events.BillOpts) error {
	return nil
}

func (w whmcsManager) Invoice(o *events.BillOpts) error {
	return nil
}

func (w whmcsManager) Nuke(o *events.BillOpts) error {
	return nil
}

func (w whmcsManager) Suspend(o *events.BillOpts) error {
	return nil
}

func (w whmcsManager) Notify(o *events.BillOpts) error {
	return nil
}
