package onetest

import (
	"testing"

	"github.com/megamsys/megamd/provision/one/cluster"
	otesting "github.com/megamsys/opennebula-go/testing"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

var _ = check.Suite(&S{})

type S struct{}

func (s *S) TestNewFakeOneProvisioner(c *check.C) {
	server, err := otesting.NewServer("127.0.0.1:5555")
	c.Assert(err, check.IsNil)
	p, err := NewFakeOneProvisioner(server.URL())
	c.Assert(err, check.IsNil)
	_, err = p.storage.RetrieveNode(server.URL())
	c.Assert(err, check.IsNil)
	defer p.Destroy()
	defer server.Stop()
}

func (s *S) TestStartMultipleServersCluster(c *check.C) {
	p, err := StartMultipleServersCluster()
	c.Assert(err, check.IsNil)
	nodes, err := p.Cluster().Nodes()
	c.Assert(err, check.IsNil)
	c.Assert(nodes, check.HasLen, 1)
	c.Assert(p.servers, check.HasLen, 1)
	defer p.Destroy()
}

func (s *S) TestDestroy(c *check.C) {
	p, err := StartMultipleServersCluster()
	c.Assert(err, check.IsNil)
	p.Destroy()
	c.Assert(p.servers, check.IsNil)
}

func (s *S) TestServers(c *check.C) {
	server, err := otesting.NewServer("127.0.0.1:0")
	c.Assert(err, check.IsNil)
	defer server.Stop()
	var p FakeOneProvisioner
	p.servers = append(p.servers, server)
	c.Assert(p.Servers(), check.DeepEquals, p.servers)
}

func (s *S) TestCluster(c *check.C) {
	var p FakeOneProvisioner
	cs, err := buildClusterStorage()
	c.Assert(err, check.IsNil)
	cluster, err := cluster.New(cs, cluster.Node{Address: "127.0.0.1:6767"})
	c.Assert(err, check.IsNil)
	p.cluster = cluster
	c.Assert(p.Cluster(), check.Equals, cluster)
}

func buildClusterStorage() (cluster.Storage, error) {
	return &cluster.MapStorage{}, nil
}
