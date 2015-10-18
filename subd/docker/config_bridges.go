package docker

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/megamsys/libgo/cmd"
)

type Bridges map[string]DockerBridge

type DockerBridge struct {
	Name    string
	Network string
	Gateway string
}

func (d DockerBridge) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 1, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Bridge", "white", "", "") + "\t" +
		cmd.Colorfy(d.Name, "blue", "", "") + "\n"))
	b.Write([]byte("Network" + "\t" + d.Network + "\n"))
	b.Write([]byte("Gateway" + "\t" + d.Gateway + "\n"))
	b.Write([]byte("---\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

func NewBridgeConfig() *Bridges {
	br := make(Bridges)
	return &br
}

func (c Bridges) String() string {
	bs := make([]string, len(c))
	for _, v := range c {
		bs = append(bs, v.String(), "\n")
	}

	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Config:", "white", "", "bold") + "\t" +
		cmd.Colorfy("bridges", "green", "", "") + "\n"))
	b.Write([]byte(strings.Join(bs, "\n")))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())

}
