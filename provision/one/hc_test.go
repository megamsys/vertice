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

package one

import (
/*	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/megamsys/libgo/hc"
	"gopkg.in/check.v1"*/
)
/*
func (s *S) TestHealthCheckDockerRegistryV2(c *check.C) {
	var request *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request = r
		w.Write([]byte("{}"))
	}))
	defer server.Close()
	if old, err := config.Get("docker:registry"); err == nil {
		defer config.Set("docker:registry", old)
	} else {
		defer config.Unset("docker:registry")
	}
	config.Set("docker:registry", server.URL+"/")
	err := healthCheckDockerRegistry()
	c.Assert(err, check.IsNil)
	c.Assert(request.URL.Path, check.Equals, "/v2/")
	c.Assert(request.Method, check.Equals, "GET")
}

func (s *S) TestHealthCheckDockerRegistryV1(c *check.C) {
	var request *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		request = r
		w.Write([]byte("pong"))
	}))
	defer server.Close()
	if old, err := config.Get("docker:registry"); err == nil {
		defer config.Set("docker:registry", old)
	} else {
		defer config.Unset("docker:registry")
	}
	config.Set("docker:registry", server.URL+"/")
	err := healthCheckDockerRegistry()
	c.Assert(err, check.IsNil)
	c.Assert(request.URL.Path, check.Equals, "/v1/_ping")
	c.Assert(request.Method, check.Equals, "GET")
}

func (s *S) TestHealthCheckDockerRegistryConfiguredWithoutScheme(c *check.C) {
	var request *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request = r
		w.Write([]byte("{}"))
	}))
	defer server.Close()
	if old, err := config.Get("docker:registry"); err == nil {
		defer config.Set("docker:registry", old)
	} else {
		defer config.Unset("docker:registry")
	}
	serverURL, _ := url.Parse(server.URL)
	config.Set("docker:registry", serverURL.Host)
	err := healthCheckDockerRegistry()
	c.Assert(err, check.IsNil)
	c.Assert(request.URL.Path, check.Equals, "/v2/")
	c.Assert(request.Method, check.Equals, "GET")
}

func (s *S) TestHealthCheckDockerRegistryFailure(c *check.C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("not pong"))
	}))
	defer server.Close()
	if old, err := config.Get("docker:registry"); err == nil {
		defer config.Set("docker:registry", old)
	} else {
		defer config.Unset("docker:registry")
	}
	config.Set("docker:registry", server.URL)
	err := healthCheckDockerRegistry()
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "unexpected status - not pong")
}

func (s *S) TestHealthCheckDockerRegistryUnconfigured(c *check.C) {
	if old, err := config.Get("docker:registry"); err == nil {
		defer config.Set("docker:registry", old)
	}
	config.Unset("docker:registry")
	err := healthCheckDockerRegistry()
	c.Assert(err, check.Equals, hc.ErrDisabledComponent)
}

func (s *S) TestHealthCheckDocker(c *check.C) {
	var request *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request = r
		w.Write([]byte("OK"))
	}))
	defer server.Close()
	var err error
	mainDockerProvisioner.cluster, err = cluster.New(nil, &cluster.MapStorage{}, cluster.Node{Address: server.URL})
	c.Assert(err, check.IsNil)
	err = healthCheckDocker()
	c.Assert(err, check.IsNil)
	c.Assert(request.Method, check.Equals, "GET")
	c.Assert(request.URL.Path, check.Equals, "/_ping")
}

func (s *S) TestHealthCheckDockerMultipleNodes(c *check.C) {
	var request *http.Request
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request = r
		w.Write([]byte("OK"))
	}))
	defer server1.Close()
	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request = r
		w.Write([]byte("OK"))
	}))
	defer server2.Close()
	var err error
	mainDockerProvisioner.cluster, err = cluster.New(nil, &cluster.MapStorage{},
		cluster.Node{Address: server1.URL}, cluster.Node{Address: server2.URL})
	c.Assert(err, check.IsNil)
	err = healthCheckDocker()
	c.Assert(err, check.Equals, hc.ErrDisabledComponent)
	c.Assert(request, check.IsNil)
}

func (s *S) TestHealthCheckDockerNoNodes(c *check.C) {
	var err error
	mainDockerProvisioner.cluster, err = cluster.New(nil, &cluster.MapStorage{})
	c.Assert(err, check.IsNil)
	err = healthCheckDocker()
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "error - no nodes available for running containers")
}

func (s *S) TestHealthCheckDockerFailure(c *check.C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong"))
	}))
	defer server.Close()
	var err error
	mainDockerProvisioner.cluster, err = cluster.New(nil, &cluster.MapStorage{}, cluster.Node{Address: server.URL})
	c.Assert(err, check.IsNil)
	err = healthCheckDocker()
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "ping failed - API error (500): something went wrong")
}
*/
