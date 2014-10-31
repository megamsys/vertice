package iaas

import (
	"gopkg.in/check.v1"

	"testing"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type S struct{}


var _ = check.Suite(&S{})

func (s *S) TestRegisterIaasProvider(c *check.C) {
	_, provider, _ := GetIaasProvider("abc")
	c.Assert(provider, check.IsNil)
}
