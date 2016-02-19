package deployd

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/opennebula-go/api"
	"github.com/megamsys/vertice/provision"
)

const (
	// Default provisioning provider for vms is OpenNebula.
	// This is just an endpoint for Megam. We could have openstack, chef, salt, puppet etc.
	DefaultProvider = "one"
	DefaultImage    = "megam"
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
	Image       string `toml:"image"`
}

func NewConfig() *Config {
	return &Config{
		Provider:    DefaultProvider,
		OneEndPoint: DefaultOneEndpoint,
		OneUserid:   DefaultOneUserid,
		OnePassword: DefaultOnePassword,
		OneTemplate: DefaultOneTemplate,
		OneZone:     DefaultOneZone,
		Certificate: "/var/lib/megam/vertice/id_rsa.pub",
		Image:       DefaultImage,
	}
}

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("\nConfig:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Deployd", "cyan", "", "") + "\n"))
	b.Write([]byte(provision.PROVIDER + "\t" + c.Provider + "\n"))
	b.Write([]byte(api.ENDPOINT + "\t" + c.OneEndPoint + "\n"))
	b.Write([]byte(api.USERID + "    \t" + c.OneUserid + "\n"))
	b.Write([]byte(api.TEMPLATE + "\t" + c.OneTemplate + "\n"))
	b.Write([]byte(api.IMAGE + "    \t" + c.Image + "\n"))
	b.Write([]byte(api.PASSWORD + "\t" + c.OnePassword + "\n"))
	b.Write([]byte("---\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

//convert the config to just a map.
func (c Config) toMap() map[string]string {
	m := make(map[string]string)
	m[api.ENDPOINT] = c.OneEndPoint
	m[api.USERID] = c.OneUserid
	m[api.PASSWORD] = c.OnePassword
	m[api.TEMPLATE] = c.OneTemplate
	m[api.IMAGE] = c.Image
	return m
}
