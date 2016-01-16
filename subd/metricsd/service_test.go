package metricsd

import (
	"gopkg.in/check.v1"
)

type S struct {
	service *Service
}

var _ = check.Suite(&S{})

func (s *S) SetUpSuite(c *check.C) {
	srv := NewService(nil, nil, nil)
	s.service = srv
	c.Assert(srv, check.NotNil)
}
