package dockertest

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"sync"

	"github.com/megamsys/megamd/provision/one/testing"
	"github.com/megamsys/megamd/cluster"
	"github.com/megamsys/megamd/db"
	"github.com/megamsys/megamd/provision"
	"github.com/megamsys/megamd/provision/docker/container"
)

type FakeDockerProvisioner struct {
	machines        map[string][]machine.Machine
	machinesMut     sync.Mutex
	storage         *cluster.MapStorage
	cluster         *cluster.Cluster
	pushes          []Push
	servers         []*testing.OneServer
	pushErrors      chan error
	moveErrors      chan error
	preparedErrors  chan error
	preparedResults chan []machine.Machine
	movings         []MachineMoving
}

func NewFakeDockerProvisioner(servers ...string) (*FakeDockerProvisioner, error) {
	var err error
	p := FakeDockerProvisioner{
		storage:         &cluster.MapStorage{},
		pushErrors:      make(chan error, 10),
		moveErrors:      make(chan error, 10),
		preparedErrors:  make(chan error, 10),
		preparedResults: make(chan []machine.Machine, 10),
		containers:      make(map[string][]machine.Machine),
	}
	nodes := make([]cluster.Node, len(servers))
	for i, server := range servers {
		nodes[i] = cluster.Node{Address: server}
	}
	p.cluster, err = cluster.New(nil, p.storage, nodes...)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func StartMultipleServersCluster() (*FakeDockerProvisioner, error) {
	server1, err := testing.NewServer("127.0.0.1:0", nil, nil)
	if err != nil {
		return nil, err
	}
	server2, err := testing.NewServer("localhost:0", nil, nil)
	if err != nil {
		return nil, err
	}
	otherUrl := strings.Replace(server2.URL(), "127.0.0.1", "localhost", 1)
	p, err := NewFakeDockerProvisioner(server1.URL(), otherUrl)
	if err != nil {
		return nil, err
	}
	p.servers = []*testing.DockerServer{server1, server2}
	return p, nil
}

func (p *FakeDockerProvisioner) SetAuthConfig(config docker.AuthConfiguration) {
	p.authConfig = config
}

func (p *FakeDockerProvisioner) Destroy() {
	for _, server := range p.servers {
		server.Stop()
	}
	p.servers = nil
}

func (p *FakeDockerProvisioner) Servers() []*testing.DockerServer {
	return p.servers
}

func (p *FakeDockerProvisioner) FailPush(errs ...error) {
	for _, err := range errs {
		p.pushErrors <- err
	}
}

func (p *FakeDockerProvisioner) Cluster() *cluster.Cluster {
	return p.cluster
}

func (p *FakeDockerProvisioner) Collection() *storage.Collection {
	conn, err := db.Conn()
	if err != nil {
		panic(err)
	}
	return conn.Collection("fake_docker_provisioner")
}

func (p *FakeDockerProvisioner) PushImage(name, tag string) error {
	p.pushes = append(p.pushes, Push{Name: name, Tag: tag})
	select {
	case err := <-p.pushErrors:
		return err
	default:
	}
	return nil
}

type Push struct {
	Name string
	Tag  string
}

func (p *FakeDockerProvisioner) Pushes() []Push {
	return p.pushes
}

func (p *FakeDockerProvisioner) RegistryAuthConfig() docker.AuthConfiguration {
	return p.authConfig
}

func (p *FakeDockerProvisioner) SetContainers(host string, containers []container.Container) {
	p.containersMut.Lock()
	defer p.containersMut.Unlock()
	dst := make([]container.Container, len(containers))
	for i, container := range containers {
		container.HostAddr = host
		dst[i] = container
	}
	p.containers[host] = dst
}

func (p *FakeDockerProvisioner) Containers(host string) []container.Container {
	p.containersMut.Lock()
	defer p.containersMut.Unlock()
	return p.containers[host]
}

func (p *FakeDockerProvisioner) AllContainers() []container.Container {
	p.containersMut.Lock()
	defer p.containersMut.Unlock()
	var result []container.Container
	for _, containers := range p.containers {
		for _, container := range containers {
			result = append(result, container)
		}
	}
	return result
}

func (p *FakeDockerProvisioner) Movings() []ContainerMoving {
	p.containersMut.Lock()
	defer p.containersMut.Unlock()
	return p.movings
}

func (p *FakeDockerProvisioner) FailMove(errs ...error) {
	for _, err := range errs {
		p.moveErrors <- err
	}
}

func (p *FakeDockerProvisioner) GetContainer(id string) (*container.Container, error) {
	container, _, err := p.findContainer(id)
	return &container, err
}

type StartContainersArgs struct {
	Endpoint  string
	App       provision.App
	Amount    map[string]int
	Image     string
	PullImage bool
}

// StartContainers starts the provided amount of containers in the provided
// endpoint.
//
// The amount is specified using a map of processes. The started containers
// will be both returned and stored internally.
func (p *FakeDockerProvisioner) StartContainers(args StartContainersArgs) ([]container.Container, error) {
	if args.PullImage {
		err := p.Cluster().PullImage(docker.PullImageOptions{Repository: args.Image}, p.RegistryAuthConfig(), args.Endpoint)
		if err != nil {
			return nil, err
		}
	}
	opts := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: args.Image,
		},
	}
	hostAddr := urlToHost(args.Endpoint)
	createdContainers := make([]container.Container, 0, len(args.Amount))
	for processName, amount := range args.Amount {
		opts.Config.Cmd = []string{processName}
		for i := 0; i < amount; i++ {
			_, cont, err := p.Cluster().CreateContainer(opts, args.Endpoint)
			if err != nil {
				return nil, err
			}
			createdContainers = append(createdContainers, container.Container{
				ID:            cont.ID,
				AppName:       args.App.GetName(),
				ProcessName:   processName,
				Type:          args.App.GetPlatform(),
				Status:        provision.StatusCreated.String(),
				HostAddr:      hostAddr,
				Version:       "v1",
				Image:         args.Image,
				User:          "root",
				BuildingImage: args.Image,
				Routable:      true,
			})
		}
	}
	p.containersMut.Lock()
	defer p.containersMut.Unlock()
	p.containers[hostAddr] = append(p.containers[hostAddr], createdContainers...)
	return createdContainers, nil
}
