package marketplacesd

import (
	"bytes"
	"fmt"
	"github.com/megamsys/libgo/cmd"
	constants "github.com/megamsys/libgo/utils"
	//"github.com/megamsys/opennebula-go/api"
	//"github.com/megamsys/vertice/provision/one"
	"strconv"
	"strings"
	"text/tabwriter"
)

const (
	// Default provisioning provider for vms is OpenNebula.
	// This is just an endpoint for Megam. We could have openstack, chef, salt, puppet etc.
	DefaultProvider = "one"
	DefaultImage    = "megam"
	ONEZONE         = "zone"
)

type Config struct {
	Enabled   bool   `json:"enabled" toml:"enabled"`
	DiskSize  string `json:"marketplace_disk_size" toml:"marketplace_disk_size"`
	CpuCore   string `json:"marketplace_cpu" toml:"marketplace_cpu"`
	RAM       string `json:"marketplace_ram" toml:"marketplace_ram"`
	Datastore string `json:"datastore" toml:"datastore"`
}

func NewConfig() *Config {
	return &Config{
		Enabled:  false,
		DiskSize: "10240",
		CpuCore:  "2",
		RAM:      "2048",
	}
}

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("\nConfig:", "white", "", "bold") + "\t" +
		cmd.Colorfy("MarketPlacesD", "cyan", "", "") + "\n"))
	b.Write([]byte("enabled      " + "\t" + strconv.FormatBool(c.Enabled) + "\n"))
	b.Write([]byte("DiskSize       " + "\t" + c.DiskSize + " MB\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

//convert the config to just an interface.
func (c Config) toMap() map[string]string {
	em := make(map[string]string)
	em[constants.ENABLED] = strconv.FormatBool(c.Enabled)
	em[constants.CPU] = c.CpuCore
	em[constants.RAM] = c.RAM
	em[constants.STORAGE] = c.DiskSize
	return em
}
