package deployd

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megamd/provision"
	"github.com/megamsys/megamd/provision/one"
)

const (
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
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Config:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Deployd", "green", "", "") + "\n"))
	b.Write([]byte(provision.PROVIDER + "\t" + c.Provider + "\n"))
	b.Write([]byte(one.ONE_ENDPOINT + "\t" + c.OneEndPoint + "\n"))
	b.Write([]byte(one.ONE_USERID + "\t" + c.OneUserid + "\n"))
	b.Write([]byte(one.ONE_PASSWORD + "\t" + c.OnePassword))
	fmt.Fprintln(w)
	w.Flush()
	return b.String()
}

//convert the config to just a map.
func (c Config) toMap() map[string]string {
	m := make(map[string]string)
	m[one.ONE_ENDPOINT] = c.OneEndPoint
	m[one.ONE_USERID] = c.OneUserid
	m[one.ONE_PASSWORD] = c.OnePassword
	m[one.ONE_TEMPLATE] = c.OneTemplate
	return m
}
