package auth

import (
	"gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct {
}

var _ = check.Suite(&S{})

func (s *S) SetUpSuite(c *check.C) {

}

func (s *S) TearDownSuite(c *check.C) {
}
