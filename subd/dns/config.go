package dns

import (
	"github.com/megamsys/libgo/cmd"
	"strconv"
)

type Config struct {
	Enabled          bool   `toml:"enabled"`
	Route53AccessKey string `toml:"route53_access_key"`
	Route53SecretKey string `toml:"route53_secret_key"`
}

func (c Config) String() string {
	table := cmd.NewTable()
	table.AddRow(cmd.Row{cmd.Colorfy("Config:", "white", "", "bold"), cmd.Colorfy("DNS", "green", "", "")})
	table.AddRow(cmd.Row{"Enabled", strconv.FormatBool(c.Enabled)})
	table.AddRow(cmd.Row{"Accesskey", c.Route53AccessKey})
	table.AddRow(cmd.Row{"Secretkey", c.Route53SecretKey})
	return table.String()
}

func NewConfig() *Config {
	return &Config{
		Enabled: true,
	}
}
