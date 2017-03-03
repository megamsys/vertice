package metricsd

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/vertice/toml"
)

const (
	DefaultCollectInterval = 10 * time.Minute
	DefaultUnits           = "1024"
	DefaultCost            = "0.1"
)

type Config struct {
	Enabled         bool          `toml:"enabled"`
	CollectInterval toml.Duration `toml:"collect_interval"`
	Snapshots       *Snapshots    `json:"snapshots" toml:"snapshots"`
	Backups         *Backups      `json:"backups" toml:"backups"`
	Skews           *Skews        `json:"skews" toml:"skews"`
}

type Snapshots struct {
	Enabled     bool   `json:"enabled" toml:"enabled"`
	StorageUnit string `json:"storage_unit" toml:"storage_unit"`
	CostPerHour string `json:"cost_per_hour" toml:"cost_per_hour"`
}

type Backups struct {
	Enabled     bool   `json:"enabled" toml:"enabled"`
	StorageUnit string `json:"storage_unit" toml:"storage_unit"`
	CostPerHour string `json:"cost_per_hour" toml:"cost_per_hour"`
}

type Skews struct {
	Enabled         bool          `json:"enabled" toml:"enabled"`
	SoftGracePeriod toml.Duration `json:"soft_grace_period" toml:"soft_grace_period"`
	SoftLimit       string        `json:"soft_limit" toml:"soft_limit"`
	HardGracePeriod toml.Duration `json:"hard_grace_period" toml:"hard_grace_period"`
	HardLimit       string        `json:"hard_limit" toml:"hard_limit"`
}

type Action struct {
}

func NewConfig() *Config {
	return &Config{
		Enabled:         false,
		CollectInterval: toml.Duration(DefaultCollectInterval),
		Snapshots: &Snapshots{
			Enabled:     false,
			StorageUnit: DefaultUnits,
			CostPerHour: DefaultCost,
		},
		Backups: &Backups{
			Enabled:     true,
			StorageUnit: DefaultUnits,
			CostPerHour: DefaultCost,
		},
		Skews: &Skews{
			Enabled:         false,
			SoftGracePeriod: toml.Duration(60 * time.Hour),
			SoftLimit:       "-5",
			HardGracePeriod: toml.Duration(120 * time.Hour),
			HardLimit:       "-10",
		},
	}
}

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Config:", "white", "", "bold") + "\t" + cmd.Colorfy("Metricsd", "cyan", "", "") + "\n"))
	b.Write([]byte("enabled" + "\t" + strconv.FormatBool(c.Enabled) + "\n"))
	b.Write([]byte("collect_interval" + "\t" + c.CollectInterval.String() + "\n"))
	b.Write([]byte("---\n"))
	b.Write([]byte(cmd.Colorfy("\nResource Bill Config:", "white", "", "bold") + "\n" + cmd.Colorfy("Bakups", "cyan", "", "") + "\n"))
	b.Write([]byte("enabled      " + "\t" + strconv.FormatBool(c.Backups.Enabled) + "\n"))
	b.Write([]byte("Cost per hour" + "\t" + c.Backups.CostPerHour + "\n"))
	b.Write([]byte("Units" + "\t" + c.Backups.StorageUnit + "\n"))
	b.Write([]byte(cmd.Colorfy("Snapshots", "cyan", "", "") + "\n"))
	b.Write([]byte("enabled      " + "\t" + strconv.FormatBool(c.Snapshots.Enabled) + "\n"))
	b.Write([]byte("Cost per hour" + "\t" + c.Snapshots.CostPerHour + "\n"))
	b.Write([]byte("Units" + "\t" + c.Snapshots.StorageUnit + "\n"))
	b.Write([]byte(cmd.Colorfy("\n Skews config:", "white", "", "bold") + "\n" + cmd.Colorfy("Actions", "cyan", "", "") + "\n"))
	b.Write([]byte("Skews enabled" + "\t" + strconv.FormatBool(c.Skews.Enabled) + "\n"))
	b.Write([]byte("Skews soft limit cost" + "\t" + c.Skews.SoftLimit + "\n"))
	b.Write([]byte("Grace period soft" + "\t" + c.Skews.SoftGracePeriod.String() + "\n"))
	b.Write([]byte("Skews hard limit cost" + "\t" + c.Skews.HardLimit + "\n"))
	b.Write([]byte("Grace period hard" + "\t" + c.Skews.HardGracePeriod.String() + "\n"))
	b.Write([]byte("---\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}
