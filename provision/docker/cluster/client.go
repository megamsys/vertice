package cluster

import (
	"bytes"
	"encoding/json"
	"net/http"
		"github.com/megamsys/vertice/carton"
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

func (d *DockerClient) NetworkRequest(url string, port string) error {
	var ips = make(map[string][]string)
	hostip := []string{}
	hostip = []string{url}
	ips[carton.HOSTIP] = hostip
	if asm, err := carton.NewAmbly(d.CartonId); err != nil {
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
