package eventsd

import (
	"gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type S struct {
	service *Service
}

var _ = check.Suite(&S{})

func (s *S) SetUpSuite(c *check.C) {
	srv := NewService(nil, nil)
	s.service = srv
	c.Assert(srv, check.NotNil)
}
