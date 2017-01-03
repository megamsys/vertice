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
	// DefaultRBDHost ceph cluster server ip
	DefaultRBDHost = "localhost"

	// Access credentials for ceph cluster server
	DefaultUsername = "vertice"
	DefaultPassword = "vertice"

	DefaultCluster = "ceph"

	DefaultCost = "1.0"
)

// Config represents the meta configuration.
type Config struct {
	Enabled bool    `json:"enabled" toml:"enabled"`
	RbdSnap RbdSnap `json:"rbd" toml:"rbd"`
}

type RbdSnap struct {
	Enabled bool     `json:"enabled" toml:"enabled"`
	Regions []Region `json:"region" toml:"region"`
}

type Region struct {
	Enabled     bool   `json:"enabled" toml:"enabled"`
	ClusterName string `json:"rbd_region" toml:"rbd_region"`
	Host        string `json:"rbd_host" toml:"rbd_host"`
	Username    string `json:"rbd_username" toml:"rbd_username"`
	Password    string `json:"rbd_password" toml:"rbd_password"`
	PoolName    string `json:"pool_name" toml:"pool_name"`
	StorageUnit string `json:"storage_unit" toml:"storage_unit"`
	CostPerHour string `json:"cost_per_hour" toml:"cost_per_hour"`
}

func NewConfig() *Config {
	rg := make([]Region, 0)
	r := Region{
		Enabled:     false,
		Host:        DefaultRBDHost,
		Username:    DefaultUsername,
		Password:    DefaultPassword,
		PoolName:    "one",
		StorageUnit: "1024",
		CostPerHour: DefaultCost,
	}

	return &Config{
		Enabled: false,
		RbdSnap: RbdSnap{
			Enabled: false,
			Regions: append(rg, r),
		},
	}
}

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("\nConfig:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Storage", "cyan", "", "") + "\n"))
	b.Write([]byte("enabled      " + "\t" + strconv.FormatBool(c.RbdSnap.Enabled) + "\n"))
	for _, v := range c.RbdSnap.Regions {
		b.Write([]byte("Rbd Server" + "\t" + v.Host + "\n"))
		b.Write([]byte("Username " + "    \t" + v.Username + "\n"))
		b.Write([]byte("Password" + "    \t" + v.Password + "\n"))
		b.Write([]byte("Poolname" + "\t" + v.PoolName + "\n"))
		b.Write([]byte("Cost per hour" + "\t" + v.CostPerHour + "\n"))
		b.Write([]byte("---\n"))
	}
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}
