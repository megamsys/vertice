package bills

//import (
	//"github.com/megamsys/vertice/carton"
//)

func init() {
	Register("scylladb", scylladbManager{})
}

type scylladbManager struct{}

func (m scylladbManager) IsEnabled() bool {
	return true
}

func (m scylladbManager) Onboard(o *BillOpts, mi map[string]string) error {
	return nil
}

func (m scylladbManager) Deduct(o *BillOpts, mi map[string]string) error {
	//b, err := carton.NewBalances(o.AccountsId)
	b, err := NewBalances(o.AccountId, mi)
	if err != nil {
		return err
	}
	//if err = b.Deduct(&carton.BalanceOpts{
	if err = b.Deduct(&BalanceOpts{
		Id:        o.AccountId,
		Consumed:  o.Consumed,
	}, mi); err != nil {
		return err
	}
	return nil
}

func (m scylladbManager) Transaction(o *BillOpts, mi map[string]string) error {
	//bt, err := carton.NewBillTransaction(o.AccountsId)
	bt, err := NewBillTransaction(o)
	if err != nil {
		return err
	}
	//if err = bt.Transact(&carton.BillTransactionOpts{
	if err = bt.Transact(mi); err != nil {
		return err
	}
	return nil
}

func (m scylladbManager) Invoice(o *BillOpts) error {
	return nil
}

func (m scylladbManager) Nuke(o *BillOpts) error {
	return nil
}

func (m scylladbManager) Suspend(o *BillOpts) error {
	return nil
}

func (m scylladbManager) Notify(o *BillOpts) error {
	return nil
}
