package bills

import (
	"github.com/megamsys/vertice/carton"
	"github.com/megamsys/vertice/events"
)

func init() {
	events.Register("scylladb", scylladbManager{})
}

type scylladbManager struct{}

func (m scylladbManager) IsEnabled(o *events.BillOpts) bool {
	return true
}

func (m scylladbManager) Onboard(o *events.BillOpts) error {
	return nil
}

func (m scylladbManager) Deduct(o *events.BillOpts) error {
	b, err := carton.NewBalances(o.AccountsId)
	if err != nil {
		return err
	}
	if err = b.Deduct(&carton.BalanceOpts{
		Id:        o.AccountsId,
		Consumed:  o.Consumed,
		Timestamp: o.Timestamp,
	}); err != nil {
		return err
	}
	return nil
}

func (m scylladbManager) Invoice(o *events.BillOpts) error {
	return nil
}

func (m scylladbManager) Nuke(o *events.BillOpts) error {
	return nil
}

func (m scylladbManager) Suspend(o *events.BillOpts) error {
	return nil
}

func (m scylladbManager) Notify(o *events.BillOpts) error {
	return nil
}
