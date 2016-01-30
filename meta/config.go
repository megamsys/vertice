package meta

import (
	"bytes"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/megamsys/libgo/cmd"
)

const (
	// DefaultRiak is the default riak if one is not provided.
	DefaultRiak = "localhost:8087"

	// DefaultApi is the default megam gateway if one is not provided.
	DefaultApi = "http://localhost:9000/v2/"

	// DefaultNSQ is the default nsqd if its not provided.
	DefaultNSQd = "localhost:4161"

	//default user
	DefaultUser = "megam"

	MEGAM_HOME = "MEGAM_HOME"
)

// Config represents the meta configuration.
type Config struct {
	Home string   `toml:"home"`
	Dir  string   `toml:"dir"`
	Logo string   `toml:"logo"`
	Riak []string `toml:"riak"`
	NSQd []string `toml:"nsqd"`
	Api  string   `toml:"api"`
	User string   `toml:"user"`
}

var MC *Config

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Config:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Meta", "cyan", "", "") + "\n"))
	b.Write([]byte("Home      " + "\t" + c.Home + "\n"))
	b.Write([]byte("Dir       " + "\t" + c.Dir + "\n"))
	b.Write([]byte("Logo      " + "\t" + c.Logo + "\n"))
	b.Write([]byte("Riak      " + "\t" + strings.Join(c.Riak, ",") + "\n"))
	b.Write([]byte("Api       " + "\t" + c.Api + "\n"))
	b.Write([]byte("NSQd      " + "\t" + strings.Join(c.NSQd, ",") + "\n"))
	b.Write([]byte("---\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

func NewConfig() *Config {
	var homeDir string
	// By default, store logs, meta and load conf files in MEGAM_HOME directory
	if os.Getenv(MEGAM_HOME) != "" {
		homeDir = os.Getenv(MEGAM_HOME)
	} else if u, err := user.Current(); err == nil {
		homeDir = u.HomeDir
	} else {
		return nil
	}

	defaultDir := filepath.Join(homeDir, "megamd/")

	// Config represents the configuration format for the megamd.
	return &Config{
		Home: homeDir,
		Dir:  defaultDir,
		User: DefaultUser,
		Riak: []string{DefaultRiak},
		Api:  DefaultApi,
		NSQd: []string{DefaultNSQd},
	}
}

func (c *Config) MkGlobal() {
	MC = c
}
