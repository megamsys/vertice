package main

import (
	"github.com/megamsys/libgo/cmd"
	"gopkg.in/check.v1"
)

func (s *S) TestMegamStartInfo(c *check.C) {
	desc := `starts the megamd daemon.

If you use the '--dry' flag megamd will do a dry run(parse conf/jsons) and exit.

`

	expected := &cmd.Info{
		Name:    "start",
		Usage:   `start [--dry] [--config]`,
		Desc:    desc,
		MinArgs: 0,
	}
	command := StartD{}
	c.Assert(command.Info(), check.DeepEquals, expected)
}