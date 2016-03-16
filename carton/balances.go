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

/*
import (
	"strconv"
	"time"

	"github.com/megamsys/vertice/db"
	"gopkg.in/yaml.v2"
)

const (
	BALANCESBUCKET = "balances"
)

type BalanceOpts struct {
	Id        string
	Consumed  int
	Timestamp time.Time
}

type Balances struct {
	AccountsId string `json:"AccountsId"`
	Name       string `json:"name"`
	Credit     string `json:"credit"`
	CreatedAt  string
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
func NewBalances(id string) (*Balances, error) {
	b := &Balances{}
	if err := db.Fetch(BALANCESBUCKET, id, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (b *Balances) Deduct(bopts *BalanceOpts) error {
	b.CreatedAt = time.Now().Local().Format(time.RFC822)

	avail, err := strconv.Atoi(b.Credit)
	if err != nil {
		return err
	}
	b.Credit = string(avail - bopts.Consumed)

	if err := db.Store(BALANCESBUCKET, bopts.Id, b); err != nil {
		return err
	}
	return nil
}
*/
