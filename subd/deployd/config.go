package deployd

import (
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megamd/provision"
)

const (
	ONE_ENDPOINT = "one_endpoint"
	ONE_USERID   = "one_userid"
	ONE_PASSWORD = "one_password"
	ONE_TEMPLATE = "one_template"
	// Default provisioning provider for vms is OpenNebula.
	// This is just an endpoint for Megam. We could have openstack, chef, salt, puppet etc.
	DefaultProvider = "one"
	// DefaultOneEndpoint is the default address that the service binds to an IaaS (OpenNebula).
	DefaultOneEndpoint = "http://localhost:2633/RPC2"

	// DefaultUserid the default userid for the IaaS service (OpenNebula).
	DefaultOneUserid = "oneadmin"

	// DefaultOnePassword is the default password for the IaaS service (OpenNebula).
	DefaultOnePassword = "password"

	// DefaultOneTemplate is the default template for the IaaS service (OpenNebula).
	DefaultOneTemplate = "megam"

	// DefaultOneZone is the default zone for the IaaS service (OpenNebula).
	DefaultOneZone = "plano01"
)

type Config struct {
	Provider    string `toml:"provider"`
	OneEndPoint string `toml:"one_endpoint"`
	OneUserid   string `toml:"one_userid"`
	OnePassword string `toml:"one_password"`
	OneTemplate string `toml:"one_template"`
	OneZone     string `toml:"one_zone"`
	Certificate string `toml:"certificate"`
}

func NewConfig() *Config {
	return &Config{
		Provider:    DefaultProvider,
		OneEndPoint: DefaultOneEndpoint,
		OneUserid:   DefaultOneUserid,
		OnePassword: DefaultOnePassword,
		OneTemplate: DefaultOneTemplate,
		OneZone:     DefaultOneZone,
		Certificate: "/var/lib/megam/megamd/id_rsa.pub",
	}
}

func (c Config) String() string {
	table := cmd.NewTable()
	table.AddRow(cmd.Row{cmd.Colorfy("Config:", "white", "", "bold"), cmd.Colorfy("Deployd", "green", "", "")})
	table.AddRow(cmd.Row{provision.PROVIDER, c.Provider})
	table.AddRow(cmd.Row{ONE_ENDPOINT, c.OneEndPoint})
	table.AddRow(cmd.Row{ONE_USERID, c.OneUserid})
	table.AddRow(cmd.Row{ONE_PASSWORD, c.OnePassword})
	table.AddRow(cmd.Row{ONE_TEMPLATE, c.OneTemplate})
	table.AddRow(cmd.Row{"", ""})
	return table.String()
}

//convert the config to just a map.
func (c Config) toMap() map[string]string {
	m := make(map[string]string)
	m[ONE_ENDPOINT] = c.OneEndPoint
	m[ONE_USERID] = c.OneUserid
	m[ONE_PASSWORD] = c.OnePassword
	m[ONE_TEMPLATE] = c.OneTemplate
	return m
}
