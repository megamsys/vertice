package alerts

import (
	"os"

	"gopkg.in/check.v1"
)

var st = os.Getenv("NIL_SLACK_TOKEN")
var ch = "ahoy"

func (s *S) TestSlack(c *check.C) {
	if st == "" {
		c.Skip("-Slack (token) not provided")
	}
	c.Assert(len(st) > 0, check.Equals, true)
	ms := NewSlack(st, ch)
	c.Assert(ms, check.NotNil)
	err := ms.Notify("test message from megamd")
	c.Assert(err, check.IsNil)
}
