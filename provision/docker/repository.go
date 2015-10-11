package docker

import (
	"errors"

	"github.com/megamsys/megamd/provision/docker/cluster"
)

var errAmbiguousContainer error = errors.New("ambiguous container name")

func (p *dockerProvisioner) GetContainer(id string) (*container.Container, error) {
	var containers []container.Container
	//stick something, so it returns the container
	lenContainers := len(containers)
	if lenContainers == 0 {
		return nil, provision.ErrUnitNotFound
	}
	if lenContainers > 1 {
		return nil, errAmbiguousContainer
	}
	return &containers[0], nil
}

func (p *dockerProvisioner) GetContainerByName(name string) (*container.Container, error) {
	var containers []container.Container
	//stick something, so it returns the container

	lenContainers := len(containers)
	if lenContainers == 0 {
		return nil, provision.ErrUnitNotFound
	}
	if lenContainers > 1 {
		return nil, errAmbiguousContainer
	}
	return &containers[0], nil
}

func (p *dockerProvisioner) listContainersByBox(boxName string) ([]container.Container, error) {
	var list []container.Container
	return list, err
}

func (p *dockerProvisioner) listBoxsForNodes(nodes []*cluster.Node) ([]string, error) {
	//coll := p.Collection()
	nodeNames := make([]string, len(nodes))
	for i, n := range nodes {
		nodeNames[i] = urlToHost(n.Address)
	}
	var appNames []string
	return appNames, nil
}
