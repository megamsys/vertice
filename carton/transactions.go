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
package carton

import (
	"time"

	"github.com/megamsys/vertice/db"
	"gopkg.in/yaml.v2"
)

const (
	TRANSACTIONBUCKET = "billtransactions"
)

type BillTransactionOpts struct {
	Id        string
	Timestamp time.Time
}

type BillTransaction struct {
	AccountsId string `json:"AccountsId"`
	Name       string `json:"name"`
	CreatedAt  string
}

func (bt *BillTransactionOpts) String() string {
	if d, err := yaml.Marshal(bt); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

func NewBillTransaction(id string) (*BillTransaction, error) {
	bt := &BillTransaction{}
	if err := db.Fetch(TRANSACTIONBUCKET, id, bt); err != nil {
		return nil, err
	}
	return bt, nil
}

func (bt *BillTransaction) Transact(topts *BillTransactionOpts) error {
	bt.CreatedAt = time.Now().Local().Format(time.RFC822)

	if err := db.Store(TRANSACTIONBUCKET, topts.Id, bt); err != nil {
		return err
	}
	return nil
}
