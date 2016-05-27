package bills

import (
	"crypto/md5"
	b64 "encoding/base64"
	"encoding/hex"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/events/addons"
	"github.com/megamsys/libgo/pairs"
	constants "github.com/megamsys/libgo/utils"
	whmcs "github.com/megamsys/whmcs_go/whmcs"
	"strings"
	"time"
)

const (
	billerName = "whmcs"
)

func init() {
	Register(billerName, createBiller())
}

type whmcsBiller struct {
	enabled bool
	apiKey  string
	domain  string
}

func createBiller() BillProvider {
	vBiller := whmcsBiller{
		enabled: false,
		apiKey:  "",
		domain:  "",
	}
	log.Debugf("%s ready", billerName)
	return vBiller
}

func (w *whmcsBiller) String() string {
	return "WHMCS:(" + w.apiKey + "," + w.domain + ")"
}

func (w whmcsBiller) IsEnabled() bool {
	return w.enabled
}

func (w whmcsBiller) Onboard(o *BillOpts, m map[string]string) error {
	log.Debugf("User Onboarding...")

	acc, err := NewAccounts(o.AccountId, m)
	if err != nil {
		return err
	}

	sDec, _ := b64.StdEncoding.DecodeString(acc.Password)

	client := whmcs.NewClient(nil, m[constants.DOMAIN])
	a := map[string]string{
		"username":    m[constants.USERNAME],
		"password":    GetMD5Hash(m[constants.PASSWORD]),
		"firstname":   acc.FirstName,
		"lastname":    acc.FirstName,
		"email":       acc.Email,
		"address1":    "Dummy address",
		"city":        "Dummy city",
		"state":       "Dummy state",
		"postcode":    "00001",
		"country":     "IN",
		"phonenumber": "981999000",
		"password2":   string(sDec),
	}

	_, res, _ := client.Accounts.Create(a)

	err = onboardNotify(acc.Email, res.Body, m)
	return err
}

func (w whmcsBiller) Deduct(o *BillOpts, m map[string]string) error {
	add := &addons.Addons{
		ProviderName: constants.WHMCS,
		AccountId:    o.AccountId,
	}
	
	err := add.Get(m)
	if err != nil {
		return err
	}

	client := whmcs.NewClient(nil, m[constants.DOMAIN])
	a := map[string]string{
		"username":      m[constants.USERNAME],
		"password":      GetMD5Hash(m[constants.PASSWORD]),
		"clientid":      add.ProviderId,
		"description":   o.AssemblyName,
		"hours":         "1",
		"amount":        o.Consumed,
		"invoiceaction": "nextcron",
	}

	_, _, err = client.Billables.Create(a)
	return err
}

func (w whmcsBiller) Transaction(o *BillOpts, m map[string]string) error {
	return nil
}

func (w whmcsBiller) Invoice(o *BillOpts) error {
	return nil
}

func (w whmcsBiller) Nuke(o *BillOpts) error {
	return nil
}

func (w whmcsBiller) Suspend(o *BillOpts) error {
	return nil
}

func (w whmcsBiller) Notify(o *BillOpts) error {
	return nil
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func onboardNotify(email string, r string, m map[string]string) error {
	cid := getClientId(r)
	if cid != "0" {
		return recordStatus(email, cid, "onboarded", m)
	}
	return nil
}

func getClientId(body string) string {
	id := "0"
	result := strings.Split(body, ";")
	for i := range result {
		if len(result[i]) > 0 {
			k := strings.Split(result[i], "=")
			if k[0] == "clientid" {
				id = k[1]
			}
		}
	}
	return id
}

func recordStatus(email, cid, status string, mi map[string]string) error {

	js := make(pairs.JsonPairs, 0)
	m := make(map[string][]string, 2)
	m["status"] = []string{status}
	js.NukeAndSet(m) //just nuke the matching output key:

	addon := &addons.Addons{
		Id:           "",
		ProviderName: constants.WHMCS,
		ProviderId:   cid,
		AccountId:    email,
		Options:      js.ToString(),
		CreatedAt:    time.Now().String(),
	}
	return addon.Onboard(mi)
}
