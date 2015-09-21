package api

import (
	"net/http"

	"gopkg.in/check.v1"
)

type IndexSuite struct{}

var _ = check.Suite(IndexSuite{})

func (IndexSuite) TestIndex(c *check.C) {
	_, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, check.IsNil)
}
