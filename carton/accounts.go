package carton

import (
	"encoding/json"
	"github.com/megamsys/libgo/api"
	"github.com/megamsys/vertice/meta"
)

const (
	ACCOUNTSBUCKET = "accounts"
)

type AccountApi struct {
	JsonClaz string  `json:"json_claz"`
	Results  Account `json:"results"`
}

type AccountsApi struct {
	JsonClaz string  `json:"json_claz"`
	Results  []Account `json:"results"`
}

type Account struct {
	Id           string   `json:"id" cql:"id"`
	Name         Name     `json:"name" cql:"name"`
	Phone        Phone    `json:"phone" cql:"phone"`
	Email        string   `json:"email" cql:"email"`
	Dates        Dates    `json:"dates" cql:"dates"`
	ApiKey       string   `json:"api_key" cql:"api_key"`
	Password     Password `json:"password" cql:"password"`
	Approval     Approval `json:"approval" cql:"approval"`
	Suspend      Suspend  `json:"suspend" cql:"suspend"`
	RegIpAddress string   `json:"registration_ip_address" cql:"registration_ip_address"`
	States       States   `json:"states" cql:"states"`
}

type Name struct {
	FirstName string `json:"first_name" cql:"first_name"`
	LastName  string `json:"last_name" cql:"last_name"`
}

type Password struct {
	Password            string `json:"password" cql:"password"`
	PasswordResetKey    string `json:"password_reset_key" cql:"password_reset_key"`
	PasswordResetSentAt string `json:"password_reset_sent_at" cql:"password_reset_sent_at"`
}

type Phone struct {
	Phone         string `json:"phone" cql:"phone"`
	PhoneVerified string `json:"phone_verified" cql:"phone_verified"`
}

type Approval struct {
	Approved     string    `json:"approved" cql:"approved"`
	ApprovedById string    `json:"approved_by_id" cql:"approved_by_id"`
	ApprovedAt   string `json:"approved_at" cql:"approved_at"`
}

type Suspend struct {
	Suspended     string    `json:"suspended" cql:"suspended"`
	SuspendedAt   string `json:"suspended_at" cql:"suspended_at"`
	SuspendedTill string `json:"suspended_till" cql:"suspended_till"`
}

type Dates struct {
	CreatedAt       string `json:"created_at" cql:"created_at"`
	LastPostedAt    string `json:"last_posted_at" cql:"last_posted_at"`
	LastEmailedAt   string `json:"last_emailed_at" cql:"last_emailed_at"`
	PreviousVisitAt string `json:"previous_visit_at" cql:"previous_visit_at"`
	FirstSeenAt     string `json:"first_seen_at" cql:"first_seen_at"`
}

type States struct {
	Authority string `json:"authority" cql:"authority"`
	Active    string `json:"active" cql:"active"`
	Blocked   string `json:"blocked" cql:"blocked"`
	Staged    string `json:"staged" cql:"staged"`
}

func NewAccounts(email string) (*Account, error) {
	return new(Account).get(newArgs(email, ""))
}

func (a *Account) get(args api.ApiArgs) (*Account, error) {
	cl := api.NewClient(args, "/accounts/" + args.Email)
	response, err := cl.Get()
	if err != nil {
		return nil, err
	}

	ac := &AccountApi{}
	err = json.Unmarshal(response, ac)
	if err != nil {
		return nil, err
	}
	return &ac.Results, nil
}

func (a *Account) GetUsers() ([]Account,  error) {
	args := newArgs(meta.MC.MasterUser, "")
	cl := api.NewClient(args, "/admin/accounts")
	response, err := cl.Get()
	if err != nil {
		return nil, err
	}

		ac := &AccountsApi{}
		err = json.Unmarshal(response, ac)
		if err != nil {
			return nil, err
		}
		return ac.Results, nil
}
