package carton

import (
	"github.com/megamsys/libgo/api"
	"github.com/megamsys/libgo/pairs"
	"encoding/json"
)

type Quota struct {
	Id          string          `json:"id" cql:"id"`
	AccountId   string          `json:"account_id" cql:"account_id"`
	Name        string          `json:"name" cql:"name"`
	JsonClaz    string          `json:"json_claz" cql:"json_claz"`
	Allowed     string          `json:"allowed" cql:"allowed"`
	AllocatedTo string          `json:"allocated_to" cql:"allocated_to"`
	CreatedAt   string          `json:"created_at" cql:"created_at"`
	UpdatedAt   string          `json:"updated_at" cql:"updated_at"`
	Inputs      pairs.JsonPairs `json:"inputs" cql:"inputs"`
}

type ApiQuota struct {
	JsonClaz string  `json:"json_claz"`
	Results  []Quota `json:"results"`
}

func (q *Quota) Update() error {
 return q.update(newArgs(q.AccountId, ""))
}

func (q *Quota) update(args api.ApiArgs) error {
	cl := api.NewClient(args, "/quota/update")
	_, err := cl.Post(q)
	if err != nil {
		return err
	}
	return nil
}

func NewQuota(accountid, id string) (*Quota, error) {
	q := new(Quota)
	q.AccountId = accountid
   return q.get(newArgs(accountid, ""),id)
}

func (q *Quota) get(args api.ApiArgs, id string) (*Quota, error) {
	cl := api.NewClient(args, "/quotas/"+id)
	response, err := cl.Get()
	if err != nil {
		return nil, err
	}
	ac := &ApiQuota{}

	//log.Debugf("Response %s :  (%s)",cmd.Colorfy("[Body]", "green", "", "bold"),string(htmlData))
	err = json.Unmarshal(response, ac)
	if err != nil {
		return nil, err
	}

	return &ac.Results[0], nil
}
