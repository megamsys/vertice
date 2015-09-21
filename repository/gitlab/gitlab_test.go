
package gitlab

import (
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}
/*
var _ = check.Suite(&GitlabSuite{})

type GitlabSuite struct {
}

func (s *GitlabSuite) SetUpSuite(c *check.C) {
	var err = error.New("testing")
	c.Assert(err, check.IsNil)
}

func (s *GitlabSuite) TearDownSuite(c *check.C) {
	var err = error.New("testing")
	c.Assert(err, check.IsNil)
}

func (s *GitlabSuite) TearDownTest(c *check.C) {
	var err = error.New("testing")
	c.Assert(err, check.IsNil)
}

func (s *GitlabSuite) TestCreateHook(c *check.C) {
	var err = error.New("testing")
	c.Assert(err, check.IsNil)
}

func (s *GitlabSuite) TestRemoveHook(c *check.C) {
	var err = error.New("testing")
	c.Assert(err, check.IsNil)
}
*/
