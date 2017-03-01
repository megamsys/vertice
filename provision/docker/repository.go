package docker

import (
	"github.com/megamsys/vertice/provision"
	"github.com/megamsys/vertice/provision/docker/container"
)

//this is essentially converting box to a container.
func (p *dockerProvisioner) GetContainerByBox(box *provision.Box) (*container.Container, error) {
	var id string
	if len(box.InstanceId) > 10 {
		id = box.InstanceId[:10]
	}
	return &container.Container{
		BoxId:     box.Id,
		CartonId:  box.CartonId,
		AccountId: box.AccountId,
		Name:      box.Name,
		BoxName:   box.GetFullName(),
		Level:     box.Level,
		Region:    box.Region,
		Status:    box.Status,
		Id:        id,
	}, nil

}

func (p *dockerProvisioner) listContainersByBox(box *provision.Box) ([]container.Container, error) {
	list := make([]container.Container, 1)
	//
	//do a query on the name to riak, and call GetContainerByBox(box)
	//

	//This is a temporary hack - sending []container.Container to assign n workers
	nx, _ := p.GetContainerByBox(box)
	list[0] = *nx
	return list, nil
}
