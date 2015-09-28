package api

import (
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct {
	token Token
}

var _ = check.Suite(&S{})

func resetHandlers() {
	megdHandlerList = []MegdHandler{}
}

func (s *S) SetUpSuite(c *check.C) {
	s.token = getTok()
	c.Assert(s.token.GetUserName(), check.Equals, "info@megam.io")
}

func getTok() Token {
	t := Token{}
	t.Token = "aaaa"
	t.UserEmail = "info@megam.io"
	return t
}
