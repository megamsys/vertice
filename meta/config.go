package meta

import (
	"bytes"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megamd/toml"
)

const (
	// DefaultHostname is the default hostname if one is not provided.
	DefaultHostname = "localhost"

	// DefaultBindAddress is the default address to bind to.
	DefaultBindAddress = ":9999"

	// DefaultRiak is the default riak if one is not provided.
	DefaultRiak = "localhost:8087"

	// DefaultApi is the default megam gateway if one is not provided.
	DefaultApi = "https://api.megam.io/v2"

	// DefaultAMQP is the default rabbitmq if one is not provided.
	DefaultAMQP = "amqp://guest:guest@localhost:5672/"

	// DefaultProvider is the default provisioner used by our engine.
	DefaultProvider = "one"

	// DefaultHeartbeatTimeout is the default heartbeat timeout for the store.
	DefaultHeartbeatTimeout = 1000 * time.Millisecond

	// DefaultElectionTimeout is the default election timeout for the store.
	DefaultElectionTimeout = 1000 * time.Millisecond

	// DefaultLeaderLeaseTimeout is the default leader lease for the store.
	DefaultLeaderLeaseTimeout = 500 * time.Millisecond
)

// Config represents the meta configuration.
type Config struct {
	Home               string        `toml:home`
	Dir                string        `toml:"dir"`
	Hostname           string        `toml:"hostname"`
	BindAddress        string        `toml:bind_address`
	Riak               []string      `toml:"riak"`
	Api                string        `toml:"api"`
	AMQP               string        `toml:"amqp"`
	Peers              []string      `toml:"-"`
	Provider           string        `toml:"provider"`
	ElectionTimeout    toml.Duration `toml:"election-timeout"`
	HeartbeatTimeout   toml.Duration `toml:"heartbeat-timeout"`
	LeaderLeaseTimeout toml.Duration `toml:"leader-lease-timeout"`
}

var MC *Config

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Config:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Meta", "green", "", "") + "\n"))
	b.Write([]byte("Home" + "\t" + c.Home + "\n"))
	b.Write([]byte("Dir" + "\t" + c.Dir + "\n"))
	b.Write([]byte("Riak" + "\t" + strings.Join(c.Riak, ",") + "\n"))
	b.Write([]byte("API" + "\t" + c.Api + "\n"))
	b.Write([]byte("AMQP" + "\t" + c.AMQP + "\n"))
	b.Write([]byte("Hostname" + "\t" + c.Hostname + "\n"))
	fmt.Fprintln(w)
	w.Flush()
	return b.String()
}

func NewConfig() *Config {
	var homeDir string
	// By default, store logs, meta and load conf files in current users home directory
	if os.Getenv("MEGAM_HOME") != "" {
		homeDir = os.Getenv("MEGAM_HOME")
	} else if u, err := user.Current(); err == nil {
		homeDir = u.HomeDir
	} else {
		return nil
		//fmt.Errorf("failed to determine home directory")
	}

	defaultDir := filepath.Join(homeDir, "megamd/meta")

	// Config represents the configuration format for the megamd.
	return &Config{
		Home:               homeDir,
		Dir:                defaultDir,
		Hostname:           DefaultHostname,
		BindAddress:        DefaultBindAddress,
		Riak:               []string{DefaultRiak},
		Api:                DefaultApi,
		AMQP:               DefaultAMQP,
		Provider:           DefaultProvider,
		ElectionTimeout:    toml.Duration(DefaultElectionTimeout),
		HeartbeatTimeout:   toml.Duration(DefaultHeartbeatTimeout),
		LeaderLeaseTimeout: toml.Duration(DefaultLeaderLeaseTimeout),
	}
}

func (c *Config) MkGlobal() {
	MC = c
}
