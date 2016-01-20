package api

import (
	"gopkg.in/check.v1"
)

func (s *S) TestGenerateUid(c *check.C) {
	c.Assert(len(Uid("")) > 0, check.Equals, true)
}
