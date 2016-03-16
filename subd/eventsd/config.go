package eventsd

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/libgo/events"
	"github.com/megamsys/libgo/events/alerts"
	"github.com/megamsys/vertice/meta"
)

type Config struct {
	Enabled bool    `toml:"enabled"`
	Mailgun Mailgun `toml:"mailgun"`
	Slack   Slack   `toml:"slack"`
	Infobip Infobip `toml:"infobip"`
	BillMgr BillMgr `toml:"bill"`
}

func NewConfig() *Config {
	return &Config{
		Enabled: false,
	}
}

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Config:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Eventsd", "cyan", "", "") + "\n"))
	b.Write([]byte("enabled      " + "\t" + strconv.FormatBool(c.Enabled) + "\n"))
	b.Write([]byte(c.Mailgun.String()))
	b.Write([]byte(c.Infobip.String()))
	b.Write([]byte(c.Slack.String() + "\n"))
	b.Write([]byte(c.BillMgr.String() + "\n"))
	b.Write([]byte("---\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

func (c Config) toMap() events.EventsConfigMap {
	em := make(events.EventsConfigMap)
	em[alerts.MAILGUN] = c.Mailgun.toMap()
	em[alerts.SLACK] = c.Slack.toMap()
	em[alerts.INFOBIP] = c.Infobip.toMap()
	em[events.BILLMGR] = c.BillMgr.toMap()
	em[alerts.META] = meta.MC.ToMap()
	return em
}
