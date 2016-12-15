package carton

import (
	ldb "github.com/megamsys/libgo/db"
	"github.com/megamsys/vertice/meta"
	"fmt"
	"encoding/json"
	"time"
)

const (
	ACCOUNTSBUCKET = "accounts"
)


type AccountDb struct {
	Id           string `json:"id" cql:"id"`
	Name         string `json:"name" cql:"name"`
	Phone        string `json:"phone" cql"phone"`
	Email        string `json:"email" cql:"email"`
	Dates        string `json:"dates" cql:"dates"`
	ApiKey       string `json:"api_key" cql:"api_key"`
	Password     string `json:"password" cql:"password"`
	Approval     string `json:"approval" cql:"approval"`
	Suspend      string `json:"suspend" cql:"suspend"`
	RegIpAddress string `json:"registration_ip_address" cql:"registration_ip_address"`
	States       string `json:"states" cql:"states"`
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
	PasswordResetSentAt time.Time `json:"password_reset_sent_at" cql:"password_reset_sent_at"`
}

type Phone struct {
	Phone         string `json:"phone" cql:"phone"`
	PhoneVerified string `json:"phone_verified" cql:"phone_verified"`
}

type Approval struct {
	Approved     string `json:"approved" cql:"approved"`
	ApprovedById string `json:"approved_by_id" cql:"approved_by_id"`
	ApprovedAt   time.Time `json:"approved_at" cql:"approved_at"`
}

type Suspend struct {
	Suspended     string `json:"suspended" cql:"suspended"`
	SuspendedAt   time.Time `json:"suspended_at" cql:"suspended_at"`
	SuspendedTill time.Time `json:"suspended_till" cql:"suspended_till"`
}

type Dates struct {
	CreatedAt       time.Time `json:"created_at" cql:"created_at"`
	LastPostedAt    time.Time `json:"last_posted_at" cql:"last_posted_at"`
	LastEmailedAt   time.Time `json:"last_emailed_at" cql:"last_emailed_at"`
	PreviousVisitAt time.Time `json:"previous_visit_at" cql:"previous_visit_at"`
	FirstSeenAt     time.Time `json:"first_seen_at" cql:"first_seen_at"`
}

type States struct {
	Authority string `json:"authority" cql:"authority"`
	Active    string `json:"active" cql:"active"`
	Blocked   string `json:"blocked" cql:"blocked"`
	Staged    string `json:"staged" cql:"staged"`
}

func NewAccounts(email string) (*Account, error) {
	a := new(AccountDb)
	a.Email = email
	return a.get()
}

func (a *AccountDb) get() (*Account, error) {
	ops := ldb.Options{
		TableName:   ACCOUNTSBUCKET,
		Pks:         []string{},
		Ccms:        []string{"email"},
		Hosts:       meta.MC.Scylla,
		Keyspace:    meta.MC.ScyllaKeyspace,
		Username:    meta.MC.ScyllaUsername,
		Password:    meta.MC.ScyllaPassword,
		PksClauses:  make(map[string]interface{}),
		CcmsClauses: map[string]interface{}{"email": a.Email},
	}
	if err := ldb.Fetchdb(ops, a); err != nil {
		return nil, err
	}
	return a.convertAccount()
}

func (a *AccountDb) convertAccount() (*Account, error) {
	b := &Account{}
	a.parseStringToStruct([]byte(a.Name), &b.Name)
	a.parseStringToStruct([]byte(a.Phone), &b.Phone)
	a.parseStringToStruct([]byte(a.Password), &b.Password)
	a.parseStringToStruct([]byte(a.Suspend), &b.Suspend)
	a.parseStringToStruct([]byte(a.Approval), &b.Approval)
	a.parseStringToStruct([]byte(a.States), &b.States)
	a.parseStringToStruct([]byte(a.Dates), &b.Dates)
	b.Id = a.Id
	b.Email = a.Email
	b.ApiKey = a.ApiKey
	b.RegIpAddress = a.RegIpAddress
	return b, nil
}

func (a *AccountDb) parseStringToStruct(b []byte, i interface{}) {
	err := json.Unmarshal(b, i)
	if err != nil {
		fmt.Println(err)
	}
}
