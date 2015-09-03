package httpd

type Config struct {
	Enabled     bool   `toml:"enabled"`
	BindAddress string `toml:"bind_address"`
}

func (c Config) String() string {
	table := NewTable()
	table.AddRow(Row{Colorfy("Config:", "white", "", "bold"), Colorfy("Httpd", "green", "", "")})
	table.AddRow(Row{"Enabled", c.Enabled})
	table.AddRow(Row{"BindAddress", c.BindAddress})
	table.AddRow(Row{"", ""})
	return table.String()
}

func NewConfig() Config {
	return Config{
		Enabled:     true,
		BindAddress: "localhost:7777",
	}
}
