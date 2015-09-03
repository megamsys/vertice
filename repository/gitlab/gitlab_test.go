// Copyright 2015 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitlab

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/megamsys/megamd/repository"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

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
