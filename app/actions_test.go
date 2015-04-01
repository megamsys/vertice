/* 
** Copyright [2013-2015] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
*/
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
