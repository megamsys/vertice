package dockertest

/*
import (
	"net"
	"net/url"
	"strings"
	"sync"

	"github.com/fsouza/go-dockerclient"
	"github.com/fsouza/go-dockerclient/testing"
	"github.com/megamsys/megamd/provision"
	"github.com/megamsys/megamd/provision/docker/cluster"
	"github.com/megamsys/megamd/provision/docker/container"
)

type ContainerMoving struct {
	ContainerID string
	HostFrom    string
	HostTo      string
}

type FakeDockerProvisioner struct {
	containers      map[string][]container.Container
	containersMut   sync.Mutex
	storage         *cluster.MapStorage
	cluster         *cluster.Cluster
	authConfig      docker.AuthConfiguration
	pushes          []Push
	servers         []*testing.DockerServer
	pushErrors      chan error
	moveErrors      chan error
	preparedErrors  chan error
	preparedResults chan []container.Container
}

func NewFakeDockerProvisioner(servers ...string) (*FakeDockerProvisioner, error) {
	var err error
	p := FakeDockerProvisioner{
		storage:         &cluster.MapStorage{},
		pushErrors:      make(chan error, 10),
		preparedErrors:  make(chan error, 10),
		preparedResults: make(chan []container.Container, 10),
		containers:      make(map[string][]container.Container),
	}
	nodes := make([]cluster.Node, len(servers))
	for i, server := range servers {
		nodes[i] = cluster.Node{Address: server}
	}
	p.cluster, err = cluster.New(p.storage, nodes...)
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

type StartContainersArgs struct {
	Endpoint  string
	Box       provision.Box
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
			_, cont, err := p.Cluster().CreateContainer(opts)
			if err != nil {
				return nil, err
			}
			createdContainers = append(createdContainers, container.Container{
				Id:            cont.ID,
				BoxName:       args.Box.GetFullName(),
				Status:        provision.StatusCreating,
				HostAddr:      hostAddr,
				Image:         args.Image,
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

func (p *FakeDockerProvisioner) findContainer(id string) (container.Container, int, error) {
	for _, containers := range p.containers {
		for i, container := range containers {
			if container.Id == id {
				return container, i, nil
			}
		}
	}
	return container.Container{}, -1, &docker.NoSuchContainer{ID: id}
}

func urlToHost(urlStr string) string {
	url, _ := url.Parse(urlStr)
	if url == nil || url.Host == "" {
		return urlStr
	}
	host, _, _ := net.SplitHostPort(url.Host)
	if host == "" {
		return url.Host
	}
	return host
}
*/
