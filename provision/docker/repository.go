package docker

import (
	"github.com/megamsys/megamd/provision"
	"github.com/megamsys/megamd/provision/docker/container"
)

//this is essentially converting box to a container.
func (p *dockerProvisioner) GetContainerByBox(box *provision.Box) (*container.Container, error) {
	return &container.Container{
		BoxId:    box.Id,
		CartonId: box.CartonId,
		Name:     box.Name,
		BoxName:  box.GetFullName(),
		Level:    box.Level,
		Status:   box.Status,
	}, nil

}

func (p *dockerProvisioner) listContainersByBox(box *provision.Box) ([]container.Container, error) {
	list := make([]container.Container, 1)
	//
	//do a query on the name to riak, and call GetContainerByBox(box)
	//


	//This is a temporary hack - sending []container.Container to assign n workers
	nx,_ := p.GetContainerByBox(box)
  list[0]	 =  *nx
	return list, nil
}
