package eventsd

import (
	"bytes"
	"fmt"
	"github.com/megamsys/libgo/cmd"
	constants "github.com/megamsys/libgo/utils"
	"strconv"
	"strings"
	"text/tabwriter"
)

type Mailer struct {
	Enabled bool   `toml:"enabled"`
	Domain  string `toml:"domain"`
	Username  string `toml:"username"`
	Password  string `toml:"password"`
	Identity  string `toml:"identity"`
	Sender  string `toml:"sender"`
	Logo    string `toml:"logo"`
	Nilavu  string `toml:"nilavu"`
}

func NewMailer() Mailer {
	return Mailer{
		Domain: "smtp.mailgun.org",
		Username: "postmaster.megambox.com",
		Password: "donotallow",
		Sender: "Megam Systems <support@megam.io>",
		Logo:   "https://s3-ap-southeast-1.amazonaws.com/megampub/images/mailers/megam_vertice.png",
		Nilavu: "localhost:3000",
	}
}

func (m Mailer) toMap() map[string]string {
	mp := make(map[string]string)
	mp[constants.ENABLED] = strconv.FormatBool(m.Enabled)
	mp[constants.USERNAME] = m.Username
	mp[constants.PASSWORD] = m.Password
	mp[constants.IDENTITY] = m.Identity
	mp[constants.DOMAIN] = m.Domain
	mp[constants.SENDER] = m.Sender
	mp[constants.NILAVU] = m.Nilavu
	mp[constants.LOGO] = m.Logo
	return mp
}

func (m Mailer) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 1, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("SMTP Mailer", "green", "", "") + "\n"))
	b.Write([]byte("enabled      " + "\t" + strconv.FormatBool(m.Enabled) + "\n"))
	b.Write([]byte("domain" + "\t" + m.Domain + "\n"))
	b.Write([]byte("sender" + "\t" + m.Sender + "\n"))
	b.Write([]byte("username" + "\t" + m.Username + "\n"))
	b.Write([]byte("password" + "\t" + m.Password + "\n"))
	b.Write([]byte("identity" + "\t" + m.Identity + "\n"))
	b.Write([]byte("nilavu    " + "\t" + m.Nilavu + "\n"))
	b.Write([]byte("logo      " + "\t" + m.Logo + "\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}




type Slack struct {
	Enabled bool   `toml:"enabled"`
	Token   string `toml:"token"`
	Channel string `toml:"channel"`
}

func (s Slack) toMap() map[string]string {
	mp := make(map[string]string)
	mp[constants.ENABLED] = strconv.FormatBool(s.Enabled)
	mp[constants.TOKEN] = s.Token
	mp[constants.CHANNEL] = s.Channel
	return mp
}

func (s Slack) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 1, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("\nSlack", "green", "", "") + "\n"))
	b.Write([]byte("enabled      " + "\t" + strconv.FormatBool(s.Enabled) + "\n"))
	b.Write([]byte("token" + "\t" + s.Token + "\n"))
	b.Write([]byte("channel" + "\t" + s.Channel + "\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

type Infobip struct {
	Enabled       bool   `toml:"enabled"`
	Username      string `toml:"username"`
	Password      string `toml:"password"`
	ApiKey        string `toml:"api_key"`
	ApplicationId string `toml:"application_id"`
	MessageId     string `toml:"message_id"`
}

func (i Infobip) toMap() map[string]string {
	mp := make(map[string]string)
	mp[constants.ENABLED] = strconv.FormatBool(i.Enabled)
	mp[constants.USERNAME] = i.Username
	mp[constants.PASSWORD] = i.Password
	mp[constants.API_KEY] = i.ApiKey
	mp[constants.APPLICATION_ID] = i.ApplicationId
	mp[constants.MESSAGE_ID] = i.MessageId
	return mp
}

func (i Infobip) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 1, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("\nInfobip", "green", "", "") + "\n"))
	b.Write([]byte("enabled      " + "\t" + strconv.FormatBool(i.Enabled) + "\n"))
	b.Write([]byte("username" + "\t" + i.Username + "\n"))
	b.Write([]byte("password" + "\t" + i.Password + "\n"))
	b.Write([]byte("api_key" + "\t" + i.ApiKey + "\n"))
	b.Write([]byte("application_id" + "\t" + i.ApplicationId + "\n"))
	b.Write([]byte("message_id" + "\t" + i.MessageId + "\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

type BillMgr struct {
	Enabled        bool     `toml:"enabled"`
	WHMCSAccessKey string   `toml:"whmcs_key"`
	WHMCSUserName  string   `toml:"whmcs_username"`
	WHMCSPassword  string   `toml:"whmcs_password"`
	WHMCSDomain    string   `toml:"whmcs_domain"`
	PiggyBanks     []string `toml:"piggybanks"`
}

func (l BillMgr) toMap() map[string]string {
	mp := make(map[string]string)
	mp[constants.ENABLED] = strconv.FormatBool(l.Enabled)
	mp[constants.USERNAME] = l.WHMCSUserName
	mp[constants.WHMCS_PASSWORD] = l.WHMCSPassword
	mp[constants.WHMCS_APIKEY] = l.WHMCSAccessKey
	mp[constants.DOMAIN] = l.WHMCSDomain
	mp[constants.PIGGYBANKS] = strings.Join(l.PiggyBanks, ",")
	return mp

}

func (l BillMgr) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 1, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("\nBillmgr", "green", "", "") + "\n"))
	b.Write([]byte("enabled      " + "\t" + strconv.FormatBool(l.Enabled) + "\n"))
	b.Write([]byte("piggybanks    " + "\t" + strings.Join(l.PiggyBanks, ",") + "\n"))
	b.Write([]byte("whmcs_username" + "\t" + l.WHMCSUserName + "\n"))
	b.Write([]byte("whmcs_password" + "\t" + l.WHMCSPassword + "\n"))
	b.Write([]byte("whmcs_key     " + "\t" + l.WHMCSAccessKey + "\n"))
	b.Write([]byte("whmcs_domain  " + "\t" + l.WHMCSDomain + "\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

type Addons struct {
}

func (l Addons) toMap() map[string]string {
	mp := make(map[string]string)
	return mp

}
