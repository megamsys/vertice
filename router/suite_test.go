package router

import (
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct {
}

var _ = check.Suite(&S{})

func (s *S) TearDownSuite(c *check.C) {
}
