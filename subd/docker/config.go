package docker

import (
	"bytes"
	"fmt"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megamd/toml"
)

const (
	DOCKER_REGISTRY  = "dock_registry"
	DOCKER_SWARM     = "dock_swarm"
	DOCKER_MEMSIZE   = "dock_mem"
	DOCKER_SWAPSIZE  = "dock_swap"
	DOCKER_CPUPERIOD = "dock_cpuperiod"
	DOCKER_CPUQUOTA  = "dock_cpuquota"

	// DefaultRegistry is the default registry for docker (public)
	DefaultRegistry = "https://hub.docker.com"

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
)

type Config struct {
	Enabled   bool          `toml:"enabled"`
	Registry  string        `toml:"registry"`
	Namespace string        `toml:"namespace"`
	Swarm     string        `toml:"swarm"`
	MemSize   int           `toml:"mem_size"`
	SwapSize  int           `toml:"swap_size"`
	CPUPeriod toml.Duration `toml:"cpu_period"`
	CPUQuota  toml.Duration `toml:"cpu_quota"`
}

func NewConfig() *Config {
	return &Config{
		Enabled:   false,
		Registry:  DefaultRegistry,
		Namespace: DefaultNamespace,
		MemSize:   DefaultMemSize,
		SwapSize:  DefaultSwapSize,
		CPUPeriod: toml.Duration(DefaultCPUPeriod),
		CPUQuota:  toml.Duration(DefaultCPUQuota),
	}
}

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Config:", "white", "", "bold") + "\t" +
		cmd.Colorfy("docker", "green", "", "") + "\n"))
	b.Write([]byte("Enabled" + "\t" + strconv.FormatBool(c.Enabled) + "\n"))
	b.Write([]byte(DOCKER_REGISTRY + "\t" + c.Registry + "\n"))
	b.Write([]byte(DOCKER_SWARM + "\t" + c.Swarm + "\n"))
	b.Write([]byte(DOCKER_MEMSIZE + "\t" + strconv.Itoa(c.MemSize) + "\n"))
	b.Write([]byte(DOCKER_SWAPSIZE + "\t" + strconv.Itoa(c.SwapSize) + "\n"))
	b.Write([]byte(DOCKER_CPUPERIOD + "\t" + c.CPUPeriod.String() + "\n"))
	b.Write([]byte(DOCKER_CPUQUOTA + "\t" + c.CPUQuota.String() + "\n"))
	fmt.Fprintln(w)
	w.Flush()
	return b.String()
}

//convert the config to just a map.
func (c Config) toMap() map[string]string {
	m := make(map[string]string)
	m[DOCKER_REGISTRY] = c.Registry
	m[DOCKER_SWARM] = c.Swarm
	m[DOCKER_MEMSIZE] = strconv.Itoa(c.MemSize)
	m[DOCKER_SWAPSIZE] = strconv.Itoa(c.SwapSize)
	m[DOCKER_CPUPERIOD] = c.CPUPeriod.String()
	m[DOCKER_CPUQUOTA] = c.CPUQuota.String()
	return m
}
