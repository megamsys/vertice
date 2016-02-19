package bills

import (
	"time"
)

const (
	SCYLLADB = "scylladb"
	WHMCS    = "whmcs"
)

var BillProviders map[string]BillProvider

//BillOpts represents a billtransaction managed by the provider
type BillOpts struct {
	AccountsId string
	Email      string
	Consumed   int
	StartTime  string
	EndTime    string
	Timestamp  time.Time
}

type BillProvider interface {
	IsEnabled() bool               //is this billing provider enabled.
	Onboard(o *BillOpts) error     //onboard an user in the billing system
	Nuke(o *BillOpts) error        //delete an user from the billing system
	Suspend(o *BillOpts) error     //suspend an user from the billing system
	Deduct(o *BillOpts) error      //deduct the balance credit
	Transaction(o *BillOpts) error //deduct the bill transaction
	Invoice(o *BillOpts) error     //invoice for a  range.
	Notify(o *BillOpts) error      //notify an user
}

// Provider returns the current configured manager, as defined in the
// configuration file.
func Provider(providerName string) BillProvider {
	if _, ok := BillProviders[providerName]; !ok {
		providerName = "nop"
	}
	return BillProviders[providerName]
}

// Register registers a new billing provider, that can be later configured
// and used.
func Register(name string, provider BillProvider) {
	if BillProviders == nil {
		BillProviders = make(map[string]BillProvider)
	}
	BillProviders[name] = provider
}
