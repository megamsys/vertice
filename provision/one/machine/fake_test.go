package machine

import (
	"github.com/megamsys/megamd/provision/one/cluster"
)

type fakeOneProvisioner struct {
	storage *cluster.MapStorage
	cluster *cluster.Cluster
}

func newFakeOneProvisioner(servers ...string) (*fakeOneProvisioner, error) {
	var err error
	p := fakeOneProvisioner{
		storage: &cluster.MapStorage{},
	}
	nodes := make([]cluster.Node, len(servers))
	for i, server := range servers {
		nodes[i] = cluster.Node{Address: server}
	}
	p.cluster, err = cluster.New(p.storage, nodes...)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *fakeOneProvisioner) Cluster() *cluster.Cluster {
	return p.cluster
}
