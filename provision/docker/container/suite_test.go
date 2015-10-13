package container

import (
	"testing"

	"github.com/fsouza/go-dockerclient"
	dtesting "github.com/fsouza/go-dockerclient/testing"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

var _ = check.Suite(&S{})

type S struct {
	p      *fakeDockerProvisioner
	server *dtesting.DockerServer
	user   string
}

func (s *S) SetUpSuite(c *check.C) {

}

func (s *S) SetUpTest(c *check.C) {
	srv, err := dtesting.NewServer("127.0.0.1:0", nil, nil)
	c.Assert(err, check.IsNil)
	s.server =  srv
	s.p, err = newFakeDockerProvisioner(s.server.URL())
	c.Assert(err, check.IsNil)
}

func (s *S) TestTearDownTest(c *check.C) {
	s.server.Stop()
}

func (s *S) removeTestContainer(c *Container) error {
	return c.Remove(s.p)
}

type newContainerOpts struct {
	BoxName     string
	Image       string
	ProcessName string
}

func (s *S) newContainer(opts newContainerOpts, p *fakeDockerProvisioner) (*Container, error) {
	if p == nil {
		p = s.p
	}
	container := Container{
		Id:          "id",
		PublicIp:          "10.10.10.10",
		HostPort:    "3333",
		HostAddr:    "127.0.0.1",
		Image:       opts.Image,
  	BoxName:     opts.BoxName,
	}
	if container.BoxName == "" {
		container.BoxName = "container"
	}

	if container.Image == "" {
		container.Image = "python:latest"
	}
	port, err := getPort()
	if err != nil {
		return nil, err
	}
	ports := map[docker.Port]struct{}{
		docker.Port(port + "/tcp"): {},
	}
	config := docker.Config{
		Image:        container.Image,
		Cmd:          []string{"ps"},
		ExposedPorts: ports,
	}
	err = p.Cluster().PullImage(docker.PullImageOptions{Repository: container.Image}, docker.AuthConfiguration{})
	if err != nil {
		return nil, err
	}
	_, c, err := p.Cluster().CreateContainer(docker.CreateContainerOptions{Config: &config})
	if err != nil {
		return nil, err
	}
	container.Id = c.ID
	return &container, nil
}
