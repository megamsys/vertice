package deployd

import (
	"bytes"
	"fmt"
	"github.com/megamsys/libgo/cmd"
	constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/opennebula-go/api"
	"github.com/megamsys/vertice/provision/one"
	"strings"
	"text/tabwriter"
	"strconv"
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

  ONEZONE  = "zone"
)

type Config struct {
	Provider    string `json:"provider" toml:"provider"`
	One one.One      `json:"one" toml:"one"`
}
/*
type deployd struct {

}
*/


func NewConfig() *Config {
	cl := make([]one.Cluster,2)
	rg := make([]one.Region,2)

  c := one.Cluster{
		ClusterId: DefaultOneCluster,
		Vnet_pri_ipv4: DefaultOneVnetPri,
		Vnet_pub_ipv4: DefaultOneVnetPub,
		Vnet_pri_ipv6: DefaultOneVnetPri,
		Vnet_pub_ipv6: DefaultOneVnetPri,
	}

	r := one.Region{
		OneZone: DefaultOneZone,
		OneEndPoint: DefaultOneEndpoint,
		OneUserid: DefaultOneUserid,
		OnePassword: DefaultOnePassword,
		OneTemplate: DefaultOneCluster,
		Certificate: "/var/lib/megam/vertice/id_rsa.pub",
		Image: DefaultImage,
		VCPUPercentage: "",
		Clusters: append(cl, c),
	}
  o := one.One{
		Enabled: true,
		Regions: append(rg, r),
		OneTemplate: DefaultOneTemplate,
		Image: DefaultImage,
		VCPUPercentage: DefaultCpuThrottle,
	}

 	return &Config{
		Provider:    DefaultProvider,
   	One: o,
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
	for i:= 0; i < len(c.One.Regions); i++ {
  b.Write([]byte(api.ONEZONE + "\t" + c.One.Regions[i].OneZone + "\n"))
	b.Write([]byte(api.ENDPOINT + "\t" + c.One.Regions[i].OneEndPoint + "\n"))
	b.Write([]byte(api.USERID + "    \t" + c.One.Regions[i].OneUserid + "\n"))
	b.Write([]byte(api.PASSWORD + "\t" + c.One.Regions[i].OnePassword + "\n"))
	b.Write([]byte(api.TEMPLATE + "\t" + c.One.Regions[i].OneTemplate + "\n"))
	b.Write([]byte(api.IMAGE + "    \t" + c.One.Regions[i].Image + "\n"))
	b.Write([]byte(api.VCPU_PERCENTAGE + "\t" + c.One.Regions[i].VCPUPercentage + "\n"))
	   for j:= 0; j < len(c.One.Regions[i].Clusters); i++ {
	fmt.Println(c.One.Regions[i].Clusters[j])
	b.Write([]byte(api.CLUSTER+ "\t" + c.One.Regions[i].Clusters[j].ClusterId + "\n"))
  }
	b.Write([]byte("---\n"))
}
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

//convert the config to just an interface.
func (c Config) toInterface() interface{} {
//	var m interface{}
	//m = c.One
  return c.One
}
