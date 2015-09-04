package docker

import (
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megamd/toml"
	"strconv"
	"strings"
	"time"
)

const (
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
	Enabled   bool     `toml:"enabled"`
	Registry  string   `toml:"registry"`
	Namespace string   `toml:"namespace"`
	Master    []string `toml:"master"`
	Bridges   map[string]string
	MemSize   int           `toml:"mem_size"`
	SwapSize  int           `toml:"swap_size"`
	CPUPeriod toml.Duration `toml:"cpu_period"`
	CPUQuota  toml.Duration `toml:"cpu_quota"`
}

func (c Config) String() string {
	bs := make([]string, len(c.Bridges))
	for k, v := range c.Bridges {
		bs = append(bs, k, v, "\n")
	}
	table := cmd.NewTable()
	table.AddRow(cmd.Row{cmd.Colorfy("Config:", "white", "", "bold"), cmd.Colorfy("Docker", "green", "", "")})
	table.AddRow(cmd.Row{"Enabled", strconv.FormatBool(c.Enabled)})
	table.AddRow(cmd.Row{"Registry", c.Registry})
	table.AddRow(cmd.Row{"Master", strings.Join(c.Master, ", ")})
	table.AddRow(cmd.Row{"Bridges", strings.Join(bs, ", ")})
	table.AddRow(cmd.Row{"MemSize", strconv.Itoa(c.MemSize)})
	table.AddRow(cmd.Row{"SwapSize", strconv.Itoa(c.SwapSize)})
	table.AddRow(cmd.Row{"CPUPeriod", c.CPUPeriod.String()})
	table.AddRow(cmd.Row{"CPUQuota", c.CPUQuota.String()})
	table.AddRow(cmd.Row{"", ""})
	return table.String()
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
