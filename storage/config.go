package storage

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
  "strconv"
	"github.com/megamsys/libgo/cmd"
)

const (
	// DefaultEndpoint rados gateway endpoint
	DefaultEndpoint = "http://localhost:7480"

  // Access credentials for radosgw
	DefaultAccessKey = "vertadmin"
	DefaultSecretKey = "vertadmin"

  DefaultZone = "in.south.tn"

  DefaultCost = "1.0"
)

// Config represents the meta configuration.
type Config struct {
	Enabled    bool    `json:"enabled" toml:"enabled"`
	RgwStorage RadosGw `json:"radosgw" toml:"radosgw"`
}

type RadosGw struct {
	Enabled bool     `json:"enabled" toml:"enabled"`
	Regions []Region `json:"region" toml:"region"`
}

type Region struct {
	Enabled     bool   `json:"enabled" toml:"enabled"`
	Zone        string `json:"radosgw_region" toml:"radosgw_region"`
	EndPoint    string `json:"radosgw_host" toml:"radosgw_host"`
	AdminAccess string `json:"admin_access_key" toml:"admin_access_key"`
	AdminSecret string `json:"admin_secret_key" toml:"admin_secret_key"`
	CostPerHour string `json:"cost_per_hour" toml:"cost_per_hour"`
}

func NewConfig() *Config {
	rg := make([]Region, 0)
	r := Region{
		Enabled:     false,
		Zone:        DefaultZone,
		EndPoint:    DefaultEndpoint,
		AdminAccess: DefaultAccessKey,
		AdminSecret: DefaultSecretKey,
		CostPerHour: DefaultCost,
	}

	return &Config{
		Enabled:    false,
		RgwStorage: RadosGw{
  		Enabled:        false,
  		Regions:        append(rg, r),
  	},
	}
}

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("\nConfig:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Storage", "cyan", "", "") + "\n"))
	b.Write([]byte("enabled      " + "\t" + strconv.FormatBool(c.RgwStorage.Enabled) + "\n"))
	for _, v := range c.RgwStorage.Regions {
		b.Write([]byte("Region" + "\t" + v.Zone + "\n"))
		b.Write([]byte("Endpoint" + "\t" + v.EndPoint + "\n"))
		b.Write([]byte("AdminAccess" + "    \t" + v.AdminAccess + "\n"))
		b.Write([]byte("AdminSecret"  + "\t" + v.AdminSecret + "\n"))
		b.Write([]byte("Cost per hour" + "\t" + v.CostPerHour + "\n"))
		b.Write([]byte("---\n"))
	}
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

//convert the config to just an interface.
func (c Config) toInterface() interface{} {
	return c.RgwStorage
}
