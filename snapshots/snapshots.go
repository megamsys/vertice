package snapshots

import (
	"bytes"
	"fmt"
	"github.com/megamsys/libgo/cmd"
	"strconv"
	"strings"
	"text/tabwriter"
)

const (
	DefaultUnits = "1024"

	DefaultCost = "1.0"
)

// Config represents the meta configuration.
type Config struct {
	Enabled bool    `json:"enabled" toml:"enabled"`
  StorageUnit string `json:"storage_unit" toml:"storage_unit"`
	CostPerHour string `json:"cost_per_hour" toml:"cost_per_hour"`
}

func NewConfig() *Config {
  return &Config{
		Enabled: false,
    StorageUnit: DefaultUnits,
    CostPerHour: DefaultCost,
	}
}

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("\nConfig:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Snapshot", "cyan", "", "") + "\n"))
  	b.Write([]byte("enabled      " + "\t" + strconv.FormatBool(c.Enabled) + "\n"))
		b.Write([]byte("Cost per hour" + "\t" + c.CostPerHour + "\n"))
    b.Write([]byte("Units" + "\t" + c.StorageUnit + "\n"))
		b.Write([]byte("---\n"))

	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}
