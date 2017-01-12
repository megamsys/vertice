package rancher
/*
import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/megamsys/libgo/hc"
)

var httpRegexp = regexp.MustCompile(`^https?://`)

func init() {
	hc.AddChecker("vertice:docker", healthCheckDocker)
}

func healthCheckDocker() (interface{}, error) {
	if !strings.Contains(mainRancherProvisioner.String(), "ready") {
		return nil, hc.ErrDisabledComponent
	}
	nodes, err := mainRancherProvisioner.Cluster().Nodes()
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
	return "Rancher Server" + nodes[0].Address + " ready", nil
}*/
