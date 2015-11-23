package docker

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/megamsys/libgo/hc"
)

var httpRegexp = regexp.MustCompile(`^https?://`)

func init() {
	hc.AddChecker("megamd:docker", healthCheckDocker)
}

func healthCheckDocker() (interface{}, error) {
	if !strings.Contains(mainDockerProvisioner.String(), "ready") {
		return nil, hc.ErrDisabledComponent
	}
	nodes, err := mainDockerProvisioner.Cluster().Nodes()
	if err != nil {
		return nil, err
	}
	if len(nodes) < 1 {
		return nil, errors.New("error - no nodes available for running containers")
	}
	client, err := nodes[0].Client()
	if err != nil {
		return nil, err
	}

	err = client.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping failed - %s", err.Error())
	}
	return "docker swarm" + nodes[0].Address + " ready", nil
}
