package rancher

import (
	"bytes"
	"fmt"
	"github.com/megamsys/libgo/cmd"
	constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/vertice/provision/rancher"
	"github.com/megamsys/vertice/provision/rancher/cluster"
	"github.com/megamsys/vertice/toml"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

const (

	// DefaultRegistry is the default registry for docker (public)
	DefaultRegistry = "https://hub.docker.com"

	DefaultProvider = "rancher"

	// DefaultNamespace is the default highlevel namespace(userid) under the registry eg: https://hub.docker.com/megam
	DefaultNamespace = "megam"

	// DefaultMemSize is the default memory size in MB used for every container launch
	DefaultMemSize = 256 * 1024 * 1024

	// DefaultSwapSize is the default memory size in MB used for every container launch
	DefaultSwapSize = 210 * 1024 * 1024

	// DefaultCPUPeriod is the default cpu period used for every container launch in ms
	DefaultCPUPeriod = 25000 * time.Millisecond

	// DefaultCPUQuota is the default cpu quota allocated for every cpu cycle for the launched container in ms
	DefaultCPUQuota = 25000 * time.Millisecond

	// DefaultOneZone is the default zone for the IaaS service.
	DefaultRancherZone = "africa"
	//DefaultOneZone is the default bridge for the docker.
	DefaultBridgeName = "eth0"

	DefaultName = "eth0"

	DefaultGulpPort = ":6666"
	DefaultNetType  = "cluster-a"

	// DefaultSwarmEndpoint is the default address that the service binds to an IaaS (Swarm).
	DefaultRancherEndpoint = "http://localhost:8080"
)

type Config struct {
	Provider string        `json:"provider" toml:"provider"`
	Rancher   rancher.Rancher `json:"rancher" toml:"rancher"`
}

func NewConfig() *Config {
	fmt.Println("***********************************************")
	rg := make([]rancher.Region, 2)
	r := rancher.Region{
		RancherZone:     DefaultRancherZone,
		RancherEndPoint:  DefaultRancherEndpoint,
		Registry:       DefaultRegistry,
		CPUPeriod:      toml.Duration(DefaultCPUPeriod),
		CPUQuota:       toml.Duration(DefaultCPUQuota),
	}

	o := rancher.Rancher{
		Enabled: true,
		Regions: append(rg, r),
		//	Namespace: DefaultNamespace,
		//MemSize:   DefaultMemSize,
		//SwapSize:  DefaultSwapSize,
	}
	return &Config{
		Provider: DefaultProvider,
		Rancher:   o,
	}
}

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Config:", "white", "", "bold") + "\t" +
		cmd.Colorfy("rancher", "cyan", "", "") + "\n"))
	b.Write([]byte(constants.PROVIDER + "\t" + c.Provider + "\n"))
	b.Write([]byte("enabled      " + "\t" + strconv.FormatBool(c.Rancher.Enabled) + "\n"))
	for  _, v := range c.Rancher.Regions {
		b.Write([]byte(cluster.RANCHER_ZONE + "\t" + v.RancherZone + "\n"))
		b.Write([]byte(cluster.RANCHER_SERVER + "\t" + v.RancherEndPoint + "\n"))
	//	b.Write([]byte(cluster.DOCKER_GULP + "\t" + c.Docker.Regions[i].DockerGulpPort + "\n"))
		//b.Write([]byte(cluster.DOCKER_REGISTRY + "\t" + c.Docker.Regions[i].Registry + "\n"))
		//	b.Write([]byte(docker.DOCKER_MEMSIZE + "       \t" + strconv.Itoa(c.Docker.Regions[i].MemSize) + "\n"))
		//	b.Write([]byte(docker.DOCKER_SWAPSIZE + "    \t" + strconv.Itoa(c.Docker.Regions[i].SwapSize) + "\n"))
		b.Write([]byte(cluster.RANCHER_CPUPERIOD + "    \t" + v.CPUPeriod.String() + "\n"))
		b.Write([]byte(cluster.RANCHER_CPUQUOTA + "    \t" + v.CPUQuota.String() + "\n"))
		b.Write([]byte("---\n"))
	}
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

func (c Config) toInterface() interface{} {
	return c.Rancher
}
