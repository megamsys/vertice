package events

import (
	"time"
)

const (
	//source flags
	SCYLLADB = "scylladb"
	WHMCS    = "whmcs"

	defaultManager = SCYLLADB
)

var providers map[string]BillProvider

/* Repository represents a repository managed by the manager. */
type BillOpts struct {
	AccountsId string
	Email      string
	Consumed   int
	BeginAudit string
	EndAudit   string
	Timestamp  time.Time
}

func (r BillOpts) GetType() string {
	return ""
}

type BillProvider interface {
	IsEnabled(o *BillOpts) bool
	Invoice(o *BillOpts) error //invoice fora  range.
	Deduct(o *BillOpts) error  //deduct the bill transaction
	Onboard(o *BillOpts) error //onboard an user in the billing system
	Nuke(o *BillOpts) error    //suspend or delete an user.
	Suspend(o *BillOpts) error
	Notify(o *BillOpts) error
}

// Provider returns the current configured manager, as defined in the
// configuration file.
func Provider(providerName string) BillProvider {
	if _, ok := providers[providerName]; !ok {
		providerName = "nop"
	}
	return providers[providerName]
}

// Register registers a new billing provider, that can be later configured
// and used.
func Register(name string, provider BillProvider) {
	if providers == nil {
		providers = make(map[string]BillProvider)
	}
	providers[name] = provider
}
