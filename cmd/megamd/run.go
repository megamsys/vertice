package main

import (
	//	"bytes"
	//	"encoding/json"
	//	"fmt"
	"github.com/megamsys/libgo/cmd"
	"launchpad.net/gnuflag"
	//	"strconv"
	//	"net/http"
)

type StartD struct {
	fs  *gnuflag.FlagSet
	dry bool
}

func (g *StartD) Info() *cmd.Info {
	desc := `starts the megamd daemon.

If you use the '--dry' flag megamd will do a dry run(parse conf/jsons) and exit.

`
	return &cmd.Info{
		Name:    "start",
		Usage:   `start [--dry] [--config]`,
		Desc:    desc,
		MinArgs: 0,
	}
}

func (c *StartD) Run(context *cmd.Context, client *cmd.Client) error {
	StartDaemon(c.dry)
	return nil
}

func (c *StartD) Flags() *gnuflag.FlagSet {
	if c.fs == nil {
		c.fs = gnuflag.NewFlagSet("megamd", gnuflag.ExitOnError)
		c.fs.BoolVar(&c.dry, "config", false, "config: the configuration file to use")
		c.fs.BoolVar(&c.dry, "c", false, "config: the configuration file to use")
		c.fs.BoolVar(&c.dry, "dry", false, "dry-run: does not start the megamd (for testing purpose)")
		c.fs.BoolVar(&c.dry, "d", false, "dry-run: does not start the megamd (for testing purpose)")
	}
	return c.fs
}
