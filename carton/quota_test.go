package carton

import (
	"gopkg.in/check.v1"
)

func (s *S) TestGetQuota(c *check.C) {
	q := new(Quota)
	q.AccountId = "vino.v@megam.io"
	q.Id = "QUO6581386565277910976"
	s.Credentials.Email = q.AccountId
  _, err := q.get(s.Credentials)
  c.Assert(err, check.IsNil)
}
