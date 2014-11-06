package main

import (
	"github.com/megamsys/libgo/cmd"
	"gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type S struct{}

var _ = check.Suite(&S{})


func (s *S) TestCommandsFromBaseManagerAreRegistered(c *check.C) {
	baseManager := cmd.BuildBaseManager("megam", version, header)
	manager := buildManager("megam")
	for name, instance := range baseManager.Commands {
		command, ok := manager.Commands[name]
		c.Assert(ok, check.Equals, true)
		c.Assert(command, check.FitsTypeOf, instance)
	}
}

/*func (s *S) TestAppStartIsRegistered(c *check.C) {
	manager := buildManager("megam")
	create, ok := manager.Commands["startapp"]
	c.Assert(ok, check.Equals, true)
	c.Assert(create, check.FitsTypeOf, &AppStart{})
}

func (s *S) TestAppStopIsRegistered(c *check.C) {
	manager := buildManager("megam")
	remove, ok := manager.Commands["stopapp"]
	c.Assert(ok, check.Equals, true)
	c.Assert(remove, check.FitsTypeOf, &AppStop{})
}
*/