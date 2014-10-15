package docker

import (
	"github.com/megamsys/megamd/provisioner"
)

func Init() {
	provisioner.Register("docker", &Docker{})
}

type Docker struct {
}

func (i *Docker) CreateCommand(assembly *provisioner.AssemblyResult) (string, error) {

	return "", nil
}
