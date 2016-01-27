package eventsd

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/megamsys/libgo/cmd"
)

type Config struct {
	Enabled bool `toml:"enabled"`
}

func NewConfig() *Config {
	return &Config{
		Enabled: false,
	}
}

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Config:", "white", "", "bold") + "\t" +
		cmd.Colorfy("events", "cyan", "", "") + "\n"))
	b.Write([]byte("enabled      " + "\t" + strconv.FormatBool(c.Enabled) + "\n"))
	b.Write([]byte("---\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}
