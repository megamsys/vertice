package docker

import (
	"net/http"
	"net/http/httptest"

	"github.com/megamsys/megamd/provision/docker/cluster"
	"gopkg.in/check.v1"
)

func (s *S) TestHealthCheckDocker(c *check.C) {
	var request *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request = r
		w.Write([]byte("OK"))
	}))
	defer server.Close()
	var err error
	mainDockerProvisioner.cluster, err = cluster.New(&cluster.MapStorage{}, cluster.Gulp{}, []cluster.Bridge{}, cluster.Node{Address: server.URL})
	c.Assert(err, check.IsNil)
	_, err = healthCheckDocker()
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
	mainDockerProvisioner.cluster, err = cluster.New(&cluster.MapStorage{}, cluster.Gulp{}, []cluster.Bridge{}, cluster.Node{Address: server1.URL}, cluster.Node{Address: server2.URL})
	c.Assert(err, check.IsNil)
	_, err = healthCheckDocker()
	c.Assert(err, check.IsNil)
	c.Assert(request, check.NotNil)
}

func (s *S) TestHealthCheckDockerNoNodes(c *check.C) {
	var err error
	mainDockerProvisioner.cluster, err = cluster.New(&cluster.MapStorage{}, cluster.Gulp{}, []cluster.Bridge{})
	c.Assert(err, check.IsNil)
	_, err = healthCheckDocker()
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
	mainDockerProvisioner.cluster, err = cluster.New(&cluster.MapStorage{}, cluster.Gulp{}, []cluster.Bridge{}, cluster.Node{Address: server.URL})
	c.Assert(err, check.IsNil)
	_, err = healthCheckDocker()
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "ping failed - API error (500): something went wrong")
}
