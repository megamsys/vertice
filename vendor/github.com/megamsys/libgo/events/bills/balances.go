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
	"strconv"
	"time"
    "strings"
	"github.com/megamsys/libgo/db"
	constants "github.com/megamsys/libgo/utils"
	"gopkg.in/yaml.v2"
)

const (
	BALANCESBUCKET = "balances"
)

type BalanceOpts struct {
	Id        string
	Consumed  string
}

type Balances struct {
	Id        string `json:"id" cql:"id"`
	AccountId string `json:"account_id" cql:"account_id"`
	Credit    string `json:"credit" cql:"credit"`
	CreatedAt string `json:"created_at" cql:"created_at"`
	UpdatedAt string `json:"updated_at" cql:"updated_at"`
}

func (b *Balances) String() string {
	if d, err := yaml.Marshal(b); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

//Temporary hack to create an assembly from its id.
//This is used by SetStatus.
//We need add a Notifier interface duck typed by Box and Carton ?
func NewBalances(id string, m map[string]string) (*Balances, error) {
	b := &Balances{}
	ops := db.Options{
		TableName:   BALANCESBUCKET,
		Pks:         []string{},
		Ccms:        []string{"account_id"},
		Hosts:       strings.Split(m[constants.SCYLLAHOST], ","),
		Keyspace:    m[constants.SCYLLAKEYSPACE],
		PksClauses:  make(map[string]interface{}),
		CcmsClauses: map[string]interface{}{"account_id": id},
	}
	if err := db.Fetchdb(ops, b); err != nil {
		return nil, err
	}
	
	return b, nil
}

func (b *Balances) Deduct(bopts *BalanceOpts, m map[string]string) error {
	avail, err := strconv.ParseFloat(b.Credit, 64)
	if err != nil {
		return err
	}
	
	consume, cerr := strconv.ParseFloat(bopts.Consumed, 64)
	if cerr != nil {
		return cerr
	}
    
    update_fields := make(map[string]interface{})
	update_fields["updated_at"] = time.Now().Local().Format(time.RFC822)
	update_fields["credit"] = strconv.FormatFloat(avail - consume, 'f', 2, 64)
	ops := db.Options{
		TableName:   BALANCESBUCKET,
		Pks:         []string{},
		Ccms:        []string{"account_id"},
		Hosts:       strings.Split(m[constants.SCYLLAHOST], ","),
		Keyspace:    m[constants.SCYLLAKEYSPACE],
		PksClauses:  make(map[string]interface{}),
		CcmsClauses: map[string]interface{}{"account_id": b.AccountId},
	}
	if err := db.Updatedb(ops, update_fields); err != nil {
		return err
	}
	
	return nil
}
