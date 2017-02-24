package marketplacesd

import (
	"bytes"
	"fmt"
	"github.com/megamsys/libgo/cmd"
	"strconv"
	"strings"
	"text/tabwriter"
)

type Config struct {
	Enabled   bool   `json:"enabled" toml:"enabled"`
}

func NewConfig() *Config {
	return &Config{	Enabled:  false}
}

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("\nConfig:", "white", "", "bold") + "\t" +
		cmd.Colorfy("MarketPlacesD", "cyan", "", "") + "\n"))
	b.Write([]byte("enabled      " + "\t" + strconv.FormatBool(c.Enabled) + "\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}
