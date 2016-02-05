package onetest

import (
	"sync"

	"github.com/megamsys/vertice/provision/one/cluster"
	"github.com/megamsys/vertice/provision/one/machine"
	"github.com/megamsys/opennebula-go/testing"
)

type FakeOneProvisioner struct {
	machines        map[string][]machine.Machine
	machinesMut     sync.Mutex
	storage         *cluster.MapStorage
	cluster         *cluster.Cluster
	servers         []*testing.OneServer
	preparedResults chan []machine.Machine
	preparedErrors  chan error
}

func NewFakeOneProvisioner(servers ...string) (*FakeOneProvisioner, error) {
	var err error
	p := FakeOneProvisioner{
		storage:         &cluster.MapStorage{},
		preparedErrors:  make(chan error, 10),
		preparedResults: make(chan []machine.Machine, 10),
		machines:        make(map[string][]machine.Machine),
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

func StartMultipleServersCluster() (*FakeOneProvisioner, error) {
	server1, err := testing.NewServer("127.0.0.1:5555")
	if err != nil {
		return nil, err
	}
	p, err := NewFakeOneProvisioner(server1.URL())
	if err != nil {
		return nil, err
	}
	p.servers = []*testing.OneServer{server1}
	return p, nil
}

func (p *FakeOneProvisioner) Destroy() {
	for _, server := range p.servers {
		server.Stop()
	}
	p.servers = nil
}

func (p *FakeOneProvisioner) Servers() []*testing.OneServer {
	return p.servers
}

func (p *FakeOneProvisioner) Cluster() *cluster.Cluster {
	return p.cluster
}

func (p *FakeOneProvisioner) Machines(host string) []machine.Machine {
	p.machinesMut.Lock()
	defer p.machinesMut.Unlock()
	return p.machines[host]
}

func (p *FakeOneProvisioner) AllContainers() []machine.Machine {
	p.machinesMut.Lock()
	defer p.machinesMut.Unlock()
	var result []machine.Machine
	for _, machines := range p.machines {
		for _, machine := range machines {
			result = append(result, machine)
		}
	}
	return result
}
