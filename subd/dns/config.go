package dns

type Config struct {
	Enabled          bool   `toml:"enabled"`
	Route53AccessKey string `toml:"route53_access_key"`
	Route53SecretKey string `toml:"route53_secret_key"`
}

func (c Config) String() string {
	table := NewTable()
	table.AddRow(Row{Colorfy("Config:", "white", "", "bold"), Colorfy("DNS", "green", "", "")})
	table.AddRow(Row{"Enabled", c.Enabled})
	table.AddRow(Row{"Accesskey", c.Route53AccessKey})
	table.AddRow(Row{"Secretkey", c.Route53SecretKey})
	return table.String()
}

func NewConfig() Config {
	return Config{
		Enabled: true,
	}
}
