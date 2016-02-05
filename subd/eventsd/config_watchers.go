package eventsd

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/vertice/events/alerts"
)

type Mailgun struct {
	ApiKey string `toml:"api_key"`
	Domain string `toml:"domain"`
}

func (m Mailgun) toMap() map[string]string {
	mp := make(map[string]string)
	mp[alerts.API_KEY] = m.ApiKey
	mp[alerts.DOMAIN] = m.Domain
	return mp
}
func (m Mailgun) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 1, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Mailgun", "green", "", "") + "\n"))
	b.Write([]byte("api_key" + "\t" + m.ApiKey + "\n"))
	b.Write([]byte("domain" + "\t" + m.Domain + "\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

type Slack struct {
	Token   string `toml:"token"`
	Channel string `toml:"channel"`
}

func (s Slack) toMap() map[string]string {
	mp := make(map[string]string)
	mp[alerts.TOKEN] = s.Token
	mp[alerts.CHANNEL] = s.Channel
	return mp
}

func (s Slack) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 1, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("\nSlack", "green", "", "") + "\n"))
	b.Write([]byte("token" + "\t" + s.Token + "\n"))
	b.Write([]byte("channel" + "\t" + s.Channel + "\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

type Infobip struct {
	Username      string `toml:"username"`
	Password      string `toml:"password"`
	ApiKey        string `toml:"api_key"`
	ApplicationId string `toml:"application_id"`
	MessageId     string `toml:"message_id"`
}

func (i Infobip) toMap() map[string]string {
	mp := make(map[string]string)
	mp[alerts.USERNAME] = i.Username
	mp[alerts.PASSWORD] = i.Password
	mp[alerts.API_KEY] = i.ApiKey
	mp[alerts.APPLICATION_ID] = i.ApplicationId
	mp[alerts.MESSAGE_ID] = i.MessageId
	return mp
}

func (i Infobip) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 1, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("\nInfobip", "green", "", "") + "\n"))
	b.Write([]byte("username" + "\t" + i.Username + "\n"))
	b.Write([]byte("password" + "\t" + i.Password + "\n"))
	b.Write([]byte("api_key" + "\t" + i.ApiKey + "\n"))
	b.Write([]byte("application_id" + "\t" + i.ApplicationId + "\n"))
	b.Write([]byte("message_id" + "\t" + i.MessageId + "\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}
