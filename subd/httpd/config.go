package httpd

import (
	"github.com/megamsys/libgo/cmd"
	"strconv"
)

type Config struct {
	Enabled     bool   `toml:"enabled"`
	BindAddress string `toml:"bind_address"`
}

func (c Config) String() string {
	table := cmd.NewTable()
	table.AddRow(cmd.Row{cmd.Colorfy("Config:", "white", "", "bold"), cmd.Colorfy("Httpd", "green", "", "")})
	table.AddRow(cmd.Row{"Enabled", strconv.FormatBool(c.Enabled)})
	table.AddRow(cmd.Row{"BindAddress", c.BindAddress})
	table.AddRow(cmd.Row{"", ""})
	return table.String()
}

func NewConfig() *Config {
	return &Config{
		Enabled:     true,
		BindAddress: "localhost:7777",
	}
}
