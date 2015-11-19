package container
/*
import (
	"github.com/megamsys/megamd/provision/docker/cluster"
)

type push struct {
	name string
	tag  string
}

type fakeDockerProvisioner struct {
	storage    *cluster.MapStorage
	cluster    *cluster.Cluster
	pushes     []push
	pushErrors chan error
}

func newFakeDockerProvisioner(servers ...string) (*fakeDockerProvisioner, error) {
	var err error
	p := fakeDockerProvisioner{
		storage:    &cluster.MapStorage{},
		pushErrors: make(chan error, 10),
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

func (p *fakeDockerProvisioner) failPush(errs ...error) {
	for _, err := range errs {
		p.pushErrors <- err
	}
}

func (p *fakeDockerProvisioner) Cluster() *cluster.Cluster {
	return p.cluster
}

func (p *fakeDockerProvisioner) PushImage(name, tag string) error {
	p.pushes = append(p.pushes, push{name: name, tag: tag})
	select {
	case err := <-p.pushErrors:
		return err
	default:
	}
	return nil
}
*/
