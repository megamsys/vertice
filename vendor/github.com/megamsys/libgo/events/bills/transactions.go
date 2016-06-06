/*
** Copyright [2013-2016] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */
package bills

import (
	"time"
    "strings"
	"github.com/megamsys/libgo/db"
	constants "github.com/megamsys/libgo/utils"
	"gopkg.in/yaml.v2"
	log "github.com/Sirupsen/logrus"
)

const (
	TRANSACTIONBUCKET = "billedhistories"
)

type BillTransactionOpts struct {
	AccountId    string
	AssemblyId   string
	AssemblyName string
	Consumed     string
}

type BillTransaction struct {
	Id            string `json:"id" cql:"id"`
	AccountId     string `json:"account_id" cql:"account_id"`
	AssemblyId    string `json:"assembly_id" cql:"assembly_id"`
	BillType      string `json:"bill_type" cql:"bill_type"`
	BillingAmount string `json:"billing_amount" cql:"billing_amount"`
	CurrencyType  string `json:"currency_type" cql:"currency_type"`
	CreatedAt     string `json:"created_at" cql:"created_at"`
}

func (bt *BillTransactionOpts) String() string {
	if d, err := yaml.Marshal(bt); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

func NewBillTransaction(topts *BillOpts) (*BillTransaction, error) {
	return &BillTransaction{
		AccountId:     topts.AccountId,
		AssemblyId:    topts.AssemblyId,
		BillType:      "VM",
		BillingAmount: topts.Consumed,
		CurrencyType:  "",
		CreatedAt:     time.Now().Local().Format(time.RFC822),
	}, nil
}

func (bt *BillTransaction) Transact(m map[string]string) error {
	ops := db.Options{
		TableName:   TRANSACTIONBUCKET,
		Pks:         []string{"bill_type", "created_at"},
		Ccms:        []string{"account_id", "assembly_id"},
		Hosts:       strings.Split(m[constants.SCYLLAHOST], ","),
		Keyspace:    m[constants.SCYLLAKEYSPACE],
		PksClauses:  map[string]interface{}{"bill_type": bt.BillType, "created_at": bt.CreatedAt},
		CcmsClauses: map[string]interface{}{"account_id": bt.AccountId, "assembly_id": bt.AssemblyId},
	}
	if err := db.Storedb(ops, bt); err != nil {
		log.Debugf(err.Error())
		return err
	}
	return nil
}
