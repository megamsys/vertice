package dns

import (
	"bytes"
	"fmt"
	"strconv"
	"text/tabwriter"

	"github.com/megamsys/libgo/cmd"
)

type Config struct {
	Enabled          bool   `toml:"enabled"`
	Route53AccessKey string `toml:"route53_access_key"`
	Route53SecretKey string `toml:"route53_secret_key"`
}

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Config:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Route", "green", "", "") + "\n"))
	b.Write([]byte("Enabled" + "\t" + strconv.FormatBool(c.Enabled) + "\n"))
	b.Write([]byte("AccessKey" + "\t" + c.Route53AccessKey + "\n"))
	b.Write([]byte("SecretKey" + "\t" + c.Route53SecretKey + "\n"))
	fmt.Fprintln(w)
	w.Flush()
	return b.String()
}

func NewConfig() *Config {
	return &Config{
		Enabled: true,
	}
}
