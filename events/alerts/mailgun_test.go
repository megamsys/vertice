package alerts

import (
	"os"

	"gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct{}

var _ = check.Suite(&S{})

var ak = os.Getenv("NIL_MAILGUN_APIKEY")

func (s *S) TestMailgun(c *check.C) {
	if ak == "" {
		c.Skip("-Mailgun (api_key) not provided")
	}
	c.Assert(len(ak) > 0, check.Equals, true)
	m := NewMailgun(ak, "ojamail.megambox.com")
	c.Assert(m, check.NotNil)
	err := m.Notify("Testing")
	c.Assert(err, check.IsNil)
}
