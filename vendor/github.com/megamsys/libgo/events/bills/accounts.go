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
	ldb "github.com/megamsys/libgo/db"
	constants "github.com/megamsys/libgo/utils"
	"strings"
)

const ACCOUNTSBUCKET = "accounts"

type Accounts struct {
	Id                  string `json:"id" cql:"id"`
	FirstName           string `json:"first_name" cql:"first_name"`
	LastName            string `json:"last_name" cql:"last_name"`
	Phone               string `json:"phone" cql"phone"`
	Email               string `json:"email" cql:"email"`
	ApiKey              string `json:"api_key" cql:"api_key"`
	Password            string `json:"password" cql:"password"`
	Authority           string `json:"authority" cql:"authority"`
	PasswordResetKey    string `json:"password_reset_key" cql:"password_reset_key"`
	PasswordResetSentAt string `json:"password_reset_sent_at" cql:"password_reset_sent_at"`
	Status              string `json:"status" cql:"status"`
}

func NewAccounts(email string, m map[string]string) (*Accounts, error) {
	a := &Accounts{}
	ops := ldb.Options{
		TableName:   ACCOUNTSBUCKET,
		Pks:         []string{},
		Ccms:        []string{"email"},
		Hosts:       strings.Split(m[constants.SCYLLAHOST], ","),
		Keyspace:    m[constants.SCYLLAKEYSPACE],
		PksClauses:  make(map[string]interface{}),
		CcmsClauses: map[string]interface{}{"email": email},
	}
	if err := ldb.Fetchdb(ops, a); err != nil {
		return nil, err
	}
	return a, nil
}
