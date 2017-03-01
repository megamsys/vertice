package cluster

import (
	"bytes"
	"encoding/json"
	"github.com/megamsys/vertice/carton"
	"net/http"
	//	"github.com/megamsys/vertice/provision"
)

const (
	DOCKER_NETWORK = "/docker/networks"
	DOCKER_LOGS    = "/docker/logs"
	HTTP           = "http://"
)

type Gulp struct {
	Port string
}

type DockerClient struct {
	ContainerName string
	ContainerId   string
	Bridge        string
	IpAddr        string
	Gateway       string
	CartonId      string
	AccountId     string
	//HostAddr       string
}

func (d *DockerClient) LogsRequest(url string, port string) error {
	url = HTTP + url + port + DOCKER_LOGS
	err := request(d, url)
	if err != nil {
		return err
	}
	return nil
}

func (d *DockerClient) NetworkRequest(url, port string) error {
	var ips = make(map[string][]string)
	hostip := []string{}
	hostip = []string{url}
	ips[carton.HOSTIP] = hostip
	if asm, err := carton.NewAssembly(d.CartonId, d.AccountId, ""); err != nil {
		return err
	} else if err = asm.NukeAndSetOutputs(ips); err != nil {
		return err
	}
	return nil
}

/*
 * Request to gulp
 */
func request(d *DockerClient, url string) error {
	res, _ := json.Marshal(&d)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(res))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, err = client.Do(req)

	if err != nil {
		return err
	}
	return nil
}
