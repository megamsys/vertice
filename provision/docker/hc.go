package docker

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/megamsys/libgo/hc"
)

var httpRegexp = regexp.MustCompile(`^https?://`)

func init() {
	hc.AddChecker("docker", healthCheckDocker)
}

func healthCheckDocker() error {
	if mainDockerProvisioner == nil {
		return hc.ErrDisabledComponent
	}
	nodes, err := mainDockerProvisioner.Cluster().Nodes()
	if err != nil {
		return err
	}
	if len(nodes) < 1 {
		return errors.New("error - no nodes available for running containers")
	}
	client, err := nodes[0].Client()
	if err != nil {
		return err
	}
	err = client.Ping()
	if err != nil {
		return fmt.Errorf("ping failed - %s", err.Error())
	}
	return nil
}
