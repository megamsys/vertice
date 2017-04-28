package carton

import (
	"encoding/json"
	"fmt"
	//log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/api"
	//  constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/vertice/meta"
	"gopkg.in/yaml.v2"
)

type Organization struct {
	Id        string `json:"id" cql:"id"`
	AccountId string `json:"accounts_id" cql:"account_id"`
	Name      string `json:"name" cql:"name"`
	JsonClaz  string `json:"json_claz" cql:"json_claz"`
	CreatedAt string `json:"created_at" cql:"created_at"`
}

type Organizations struct {
	JsonClaz string         `json:"json_claz"`
	Results  []Organization `json:"results"`
}

func (a *Organization) String() string {
	if d, err := yaml.Marshal(a); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

func NewOrg(email, id string) (*Organization, error) {
	return new(Organization).get(newArgs(email, id))
}

func OrgBox() ([]Organization, error) {
	return new(Organization).adminGets(newArgs(meta.MC.MasterUser, ""))
}

func (a *Organization) get(args api.ApiArgs) (*Organization, error) {
	cl := api.NewClient(args, "/organizations/"+args.Org_Id)
	response, err := cl.Get()
	if err != nil {
		return nil, err
	}

	ac := &Organizations{}
	err = json.Unmarshal(response, ac)
	if err != nil {
		return nil, err
	}
	return &ac.Results[0], nil
}

func (a *Organization) gets(args api.ApiArgs) ([]Organization, error) {
	cl := api.NewClient(args, "/organizations")
	response, err := cl.Get()
	if err != nil {
		return nil, err
	}
	ac := &Organizations{}
	err = json.Unmarshal(response, ac)
	if err != nil {
		return nil, err
	}
	if len(ac.Results) > 0 {
		return ac.Results, nil
	}
	return nil, fmt.Errorf("No records found")
}

func (a *Organization) adminGets(args api.ApiArgs) ([]Organization, error) {
	cl := api.NewClient(args, "/admin/organizations")
	response, err := cl.Get()
	if err != nil {
		return nil, err
	}
	ac := &Organizations{}
	err = json.Unmarshal(response, ac)
	if err != nil {
		return nil, err
	}
	if len(ac.Results) > 0 {
		return ac.Results, nil
	}
	return nil, fmt.Errorf("No records found")
}
