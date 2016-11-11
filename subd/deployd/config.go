package deployd

import (
	"bytes"
	"fmt"
	"github.com/megamsys/libgo/cmd"
	constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/opennebula-go/api"
	"github.com/megamsys/vertice/provision/one"
	"strconv"
	"strings"
	"text/tabwriter"
)

const (
	// Default provisioning provider for vms is OpenNebula.
	// This is just an endpoint for Megam. We could have openstack, chef, salt, puppet etc.
	DefaultProvider = "one"
	DefaultImage    = "megam"

	DefaultCpuThrottle = "1"
	// DefaultOneEndpoint is the default address that the service binds to an IaaS (OpenNebula).
	DefaultOneEndpoint = "http://localhost:2633/RPC2"

	// DefaultUserid the default userid for the IaaS service (OpenNebula).
	DefaultOneUserid = "oneadmin"

	// DefaultOnePassword is the default password for the IaaS service (OpenNebula).
	DefaultOnePassword = "password"

	// DefaultOneTemplate is the default template for the IaaS service (OpenNebula).
	DefaultOneTemplate = "megam"

	// DefaultOneZone is the default master zone for the IaaS service (OpenNebula).
	DefaultOneMasterZone = "OpenNebula"

	// DefaultOneZone is the default zone for the IaaS service (OpenNebula).
	DefaultOneZone = "africa"

	//DefaultOneCluster is the default cluster for Host in the Iaas service (OpenNebula)
	DefaultOneCluster = "cluster-a"

	//DefaultOneCluster is the default cluster for Host in the Iaas service (OpenNebula)
	DefaultOneVnetPri = "vnet-pri"

	//DefaultOneCluster is the default cluster for Host in the Iaas service (OpenNebula)
	DefaultOneVnetPub = "vnet-pub"

	ONEZONE = "zone"
)

type Config struct {
	Provider string  `json:"provider" toml:"provider"`
	One      one.One `json:"one" toml:"one"`
}

/*
type deployd struct {

}
*/

func NewConfig() *Config {
	cl := make([]one.Cluster, 2)
	rg := make([]one.Region, 2)

	c := one.Cluster{
		Enabled:       false,
		ClusterId:     DefaultOneCluster,
		Vnet_pri_ipv4: DefaultOneVnetPri,
		Vnet_pub_ipv4: DefaultOneVnetPub,
		Vnet_pri_ipv6: DefaultOneVnetPri,
		Vnet_pub_ipv6: DefaultOneVnetPri,
	}

	r := one.Region{
		OneZone:        DefaultOneZone,
		OneEndPoint:    DefaultOneEndpoint,
		OneUserid:      DefaultOneUserid,
		OnePassword:    DefaultOnePassword,
		OneTemplate:    DefaultOneCluster,
		MemoryUnit:     "1024",
		CpuUnit:        "1",
		DiskUnit:       "10240",
		Certificate:    "/var/lib/megam/vertice/id_rsa.pub",
		Image:          DefaultImage,
		VCPUPercentage: "",
		Clusters:       append(cl, c),
	}
	o := one.One{
		Enabled:        true,
		Regions:        append(rg, r),
		OneTemplate:    DefaultOneTemplate,
		Image:          DefaultImage,
		VCPUPercentage: DefaultCpuThrottle,
	}

	return &Config{
		Provider: DefaultProvider,
		One:      o,
	}
}

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("\nConfig:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Deployd", "cyan", "", "") + "\n"))
	b.Write([]byte(constants.PROVIDER + "\t" + c.Provider + "\n"))
	b.Write([]byte("enabled      " + "\t" + strconv.FormatBool(c.One.Enabled) + "\n"))
	for _, v := range c.One.Regions {
		b.Write([]byte(api.ONEZONE + "\t" + v.OneZone + "\n"))
		b.Write([]byte(api.ENDPOINT + "\t" + v.OneEndPoint + "\n"))
		b.Write([]byte(api.USERID + "    \t" + v.OneUserid + "\n"))
		b.Write([]byte(api.PASSWORD + "\t" + v.OnePassword + "\n"))
		b.Write([]byte(api.TEMPLATE + "\t" + v.OneTemplate + "\n"))
		b.Write([]byte(api.IMAGE + "    \t" + v.Image + "\n"))
		b.Write([]byte(api.VCPU_PERCENTAGE + "\t" + v.VCPUPercentage + "\n"))
		b.Write([]byte("Resource Basic Units :\n\t\t Memory: " + v.MemoryUnit + "\tCpu: " + v.CpuUnit + "\tStorage: " + v.DiskUnit+ "\n"))
		for _, k := range v.Clusters {
			if k.Enabled {
				b.Write([]byte(api.CLUSTER + "\t" + k.ClusterId + "\n"))
			}
		}
		b.Write([]byte("---\n"))
	}
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

//convert the config to just an interface.
func (c Config) toInterface() interface{} {
	return c.One
}
