package app

import (
/*	"errors"
	"fmt"
	"github.com/tsuru/config"
	"github.com/indykish/gulp/action"
	"http://gopkg.in/check.v1"
	"sort"
	"strings"*/
)

/*
func (s *S) TestCreateRepositoryForwardInvalidType(c *check.C) {
	ctx := action.FWContext{Params: []interface{}{"something"}}
	_, err := createRepository.Forward(ctx)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "First parameter must be App or *App.")
}

func (s *S) TestCreateRepositoryBackward(c *check.C) {
	h := testHandler{}
	ts := s.t.StartGandalfTestServer(&h)
	defer ts.Close()
	app := App{Name: "someapp"}
	ctx := action.BWContext{FWResult: &app, Params: []interface{}{app}}
	createRepository.Backward(ctx)
	c.Assert(h.url[0], check.Equals, "/repository/someapp")
	c.Assert(h.method[0], check.Equals, "DELETE")
	c.Assert(string(h.body[0]), check.Equals, "null")
}
*/
