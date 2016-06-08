package addons

import (
	log "github.com/Sirupsen/logrus"
	ldb "github.com/megamsys/libgo/db"
	constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/libgo/events/alerts"
	"strings"
	"time"
)

const (
	ADDONSBUCKET = "addons"
	PROVIDER_NAME = "provider_name"
	PROVIDER_ID = "provider_id"
	)

type Addons struct {
	Id           string   `json:"id" cql:"id"`
	ProviderName string   `json:"provider_name" cql:"provider_name"`
	ProviderId   string   `json:"provider_id" cql:"provider_id"`
	AccountId    string   `json:"account_id" cql:"account_id"`
	Options      []string `json:"options" cql:"options"`
	CreatedAt    string   `json:"created_at" cql:"created_at"`
}

func NewAddons(edata alerts.EventData) *Addons {
	return &Addons{
		Id: "",
		ProviderName: edata.M[PROVIDER_NAME],
		ProviderId: edata.M[PROVIDER_ID],
		AccountId: edata.M[constants.ACCOUNT_ID],
		Options: edata.D,
		CreatedAt: time.Now().String(),
	}
}

func (s *Addons) Onboard(m map[string]string) error {
	ops := ldb.Options{
		TableName:   ADDONSBUCKET,
		Pks:         []string{PROVIDER_NAME},
		Ccms:        []string{constants.ACCOUNT_ID},
		Hosts:       strings.Split(m[constants.SCYLLAHOST], ","),
		Keyspace:    m[constants.SCYLLAKEYSPACE],
		PksClauses:  map[string]interface{}{PROVIDER_NAME: s.ProviderName},
		CcmsClauses: map[string]interface{}{constants.ACCOUNT_ID: s.AccountId},
	}
	if err := ldb.Storedb(ops, s); err != nil {
		log.Debugf(err.Error())
		return err
	}

	return nil
}

func (s *Addons) Get(m map[string]string) error {
	ops := ldb.Options{
		TableName:   ADDONSBUCKET,
		Pks:         []string{PROVIDER_NAME},
		Ccms:        []string{constants.ACCOUNT_ID},
		Hosts:       strings.Split(m[constants.SCYLLAHOST], ","),
		Keyspace:    m[constants.SCYLLAKEYSPACE],
		PksClauses:  map[string]interface{}{PROVIDER_NAME: s.ProviderName},
		CcmsClauses: map[string]interface{}{constants.ACCOUNT_ID: s.AccountId},	
	}
	if err := ldb.Fetchdb(ops, s); err != nil {
		return err
	}
	return nil
}




