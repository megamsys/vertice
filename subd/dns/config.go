package dns

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/megamsys/libgo/cmd"
)

type Config struct {
	Enabled   bool   `toml:"enabled"`
	AccessKey string `toml:"access_key"`
	SecretKey string `toml:"secret_key"`
}

var R53 *Config

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Config:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Route", "cyan", "", "") + "\n"))
	b.Write([]byte("Enabled  " + "\t" + strconv.FormatBool(c.Enabled) + "\n"))
	b.Write([]byte("AccessKey" + "\t" + c.AccessKey + "\n"))
	b.Write([]byte("SecretKey" + "\t" + c.SecretKey + "\n"))
	b.Write([]byte("---\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

func NewConfig() *Config {
	return &Config{
		Enabled: true,
	}
}

func (c *Config) MkGlobal() {
	R53 = c
}
