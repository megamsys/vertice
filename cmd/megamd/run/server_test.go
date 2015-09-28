package run

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

import (
	"fmt"
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/megamsys/megamd/subd/deployd"
	"github.com/megamsys/megamd/subd/httpd"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type S struct {
	srv *Server
	cf  *Config
}

var _ = check.Suite(&S{})

// NewTestConfig returns the default config with temporary paths.
func NewTestConfig() *Config {
	cm := NewConfig()
	u, _ := os.Getwd()
	if _, err := toml.DecodeFile(u+"/megamd.conf", &cm); err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return cm
}

// OpenServer opens a test server.
func OpenServer(c *Config) *Server {
	s, _ := NewServer(c, "0.9.1")
	_ = s.Open()
	return s
}

func (s *S) SetUpSuite(c *check.C) {
	s.cf = NewTestConfig()
	c.Assert(s.cf, check.NotNil)
	s.srv = OpenServer(s.cf)
}

func (s *S) TearDownSuite(c *check.C) {
	if s.srv != nil {
		ss := *s.srv
		for _, service := range ss.Services {
			service.Close()
		}
	}
}

func (s *S) TestServer(c *check.C) {
	c.Assert(s.srv, check.NotNil)
	c.Assert(s.cf, check.NotNil)
	ss := *s.srv
	c.Assert(ss.Services, check.HasLen, 2)
}

// URL returns the base URL for the httpd endpoint.

func (s *S) TestRabbitMQRunning(c *check.C) {
	res := false
	ss := *s.srv
	for _, service := range ss.Services {
		if _, ok := service.(*deployd.Service); ok {
			res = true
			break
		}
	}
	c.Assert(res, check.Equals, true)
}

// URL returns the base URL for the httpd endpoint.
func (s *S) TestHttpdRunning(c *check.C) {
	res := false
	ss := *s.srv
	for _, service := range ss.Services {
		if _, ok := service.(*httpd.Service); ok {
			res = true
			break
		}
	}
	c.Assert(res, check.Equals, true)
}
