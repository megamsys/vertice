package docker

import (
	"os"
	"time"

	"github.com/megamsys/megamd/toml"
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
	Enabled   bool    		 	`toml:"enabled"`
	Registry  string   			`toml:"registry"`
	Namespace string   			`toml:"namespace"`
	Master    []string 			`toml:"master"`
	Bridges   map[string]bridge
	MemSize   int           `toml:"mem_size"`
	SwapSize  int           `toml:"swap_size"`
	CPUPeriod toml.Duration `toml:"cpu_period"`
	CPUQuota  toml.Duration `toml:"cpu_quota"`
}

func (c Config) String() string {
	table := NewTable()
	table.AddRow(Row{Colorfy("Config:", "white", "", "bold"), Colorfy("Docker", "green", "", "")})
	table.AddRow(Row{"Enabled", c.Enabled})
	table.AddRow(Row{"Registry", c.Registry})
	table.AddRow(Row{"Master", c.Master.String()})
	table.AddRow(Row{"Bridges", c.Bridges.String()})
	table.AddRow(Row{"MemSize", c.MemSize})
	table.AddRow(Row{"SwapSize", c.SwapSize})
	table.AddRow(Row{"CPUPeriod", c.CPUPeriod})
	table.AddRow(Row{"CPUQuota", c.CPUQuota})
	table.AddRow(Row{"", ""})
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
