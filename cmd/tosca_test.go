package cmd

import (
  "gopkg.in/check.v1"
)



func (s *S) TestTOSCA1(c *check.C) {
  st2 := NewTOSCA()
  c.Assert(st2, check.IsNil)

}
