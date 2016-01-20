package metrix

import (
	"io/ioutil"
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type S struct {
	testjson []byte
}

var _ = check.Suite(&S{})

func (s *S) SetUpSuite(c *check.C) {
	b, e := ioutil.ReadFile("fixtures/one.json")
	c.Assert(e, check.IsNil)
	s.testjson = b
	c.Assert(s.testjson, check.NotNil)
}
