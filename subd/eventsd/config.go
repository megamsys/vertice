package eventsd

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/libgo/events"
	constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/vertice/meta"
)

type Config struct {
	Enabled bool    `toml:"enabled"`
	Mailer  Mailer  `toml:"smtp"`
	Slack   Slack   `toml:"slack"`
	Infobip Infobip `toml:"infobip"`
	BillMgr BillMgr `toml:"bill"`
	Addons  Addons  `toml:"addons"`
}

func NewConfig() *Config {
	return &Config{
		Enabled: true,
	}
}

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Config:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Eventsd", "cyan", "", "") + "\n"))
	b.Write([]byte(c.Mailer.String()))
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
	em[constants.SMTP] = c.Mailer.toMap()
	em[constants.SLACK] = c.Slack.toMap()
	em[constants.INFOBIP] = c.Infobip.toMap()
	em[constants.BILLMGR] = c.BillMgr.toMap()
	em[constants.ADDONS] = c.Addons.toMap()
	em[constants.META] = meta.MC.ToMap()
	return em
}
