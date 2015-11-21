package cluster

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	DOCKER_NETWORK = "/docker/networks"
	DOCKER_LOGS    = "/docker/logs"
)

type Gulp struct {
	Port string
}

type DockerClient struct {
	ContainerName string
	ContainerID   string
	Bridge        string
	ContainerId   string
	IpAddr        string
	Gateway       string
}

func (d *DockerClient) LogsRequest(url string, port string) error {
	url = url + port + DOCKER_LOGS
	err := request(d, url)
	if err != nil {
		return err
	}
	return nil
}

func (d *DockerClient) NetworkRequest(url string, port string) error {
	url = url + port + DOCKER_NETWORK
	err := request(d, url)
	if err != nil {
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
