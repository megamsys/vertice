/* 
** Copyright [2013-2015] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
*/
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
