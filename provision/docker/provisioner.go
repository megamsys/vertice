package docker

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/libgo/utils"
	constants "github.com/megamsys/libgo/utils"
	lb "github.com/megamsys/vertice/logbox"
	"github.com/megamsys/vertice/provision"
	"github.com/megamsys/vertice/provision/docker/cluster"
	"github.com/megamsys/vertice/provision/docker/container"
	"github.com/megamsys/vertice/repository"
	"github.com/megamsys/vertice/router"
	_ "github.com/megamsys/vertice/router/route53"
	"github.com/megamsys/vertice/toml"
)

var mainDockerProvisioner *dockerProvisioner

func init() {
	mainDockerProvisioner = &dockerProvisioner{}
	provision.Register("docker", mainDockerProvisioner)
}

type dockerProvisioner struct {
	cluster        *cluster.Cluster
	collectionName string
	storage        cluster.Storage
}
type Docker struct {
	Enabled bool     `json:"enabled" toml:"enabled"`
	Regions []Region `json:"region" toml:"region"`
}

type Region struct {
	DockerZone     string        `json:"docker_zone" toml:"docker_zone"`
	SwarmEndPoint  string        `json:"swarm" toml:"swarm"`
	DockerGulpPort string        `json:"gulp_port" toml:"gulp_port"`
	Registry       string        `json:"registry" toml:"registry"`
	CPUPeriod      toml.Duration `json:"cpu_period" toml:"cpu_period"`
	CPUQuota       toml.Duration `json:"cpu_quota" toml:"cpu_quota"`
	CpuCostPerHour string 			 `json:"cpu_cost_per_hour" toml:"cpu_cost_per_hour"`
	RamCostPerHour string 			 `json:"ram_cost_per_hour" toml:"ram_cost_per_hour"`
	CpuUnit        string        `json:"cpu_unit" toml:"cpu_unit"`
	MemoryUnit     string        `json:"memory_unit" toml:"memory_unit"`
	DiskUnit       string        `json:"disk_unit" toml:"disk_unit"`
}

func (p *dockerProvisioner) Cluster() *cluster.Cluster {
	if p.cluster == nil {
		panic("✗ docker cluster")
	}
	return p.cluster
}

func (p *dockerProvisioner) String() string {
	if p.cluster == nil {
		return "✗ docker cluster"
	}
	return "ready"
}

func (p *dockerProvisioner) Initialize(m interface{}) error {
	return p.initDockerCluster(m)
}

func (p *dockerProvisioner) initDockerCluster(i interface{}) error {
	var err error
	if p.storage == nil {
		p.storage, err = buildClusterStorage()
		if err != nil {
			return err
		}
	}
	if w, ok := i.(Docker); ok {
		var nodes []cluster.Node
		for i := 0; i < len(w.Regions); i++ {
			m := w.Regions[i].toMap()
			n := cluster.Node{
				Address:  m[cluster.DOCKER_SWARM], //swarm endpoint
				Metadata: m,
			}
			nodes = append(nodes, n)
		}

		//register nodes using the map.
		p.cluster, err = cluster.New(p.storage, nodes...)
		if err != nil {
			return err
		}
	}
	return nil
}

//convert the config to just a map.

func (c Region) toMap() map[string]string {
	m := make(map[string]string)
	m[cluster.DOCKER_ZONE] = c.DockerZone
	m[cluster.DOCKER_SWARM] = c.SwarmEndPoint
	m[cluster.DOCKER_GULP] = c.DockerGulpPort
	m[cluster.DOCKER_REGISTRY] = c.Registry
	m[cluster.DOCKER_CPUPERIOD] = c.CPUPeriod.String()
	m[cluster.DOCKER_CPUQUOTA] = c.CPUQuota.String()
	return m
}

func buildClusterStorage() (cluster.Storage, error) {
	return &cluster.MapStorage{}, nil
}

func getRouterForBox(box *provision.Box) (router.Router, error) {
	routerName, err := box.GetRouter()
	if err != nil {
		return nil, err
	}
	return router.Get(routerName)
}

func (p *dockerProvisioner) StartupMessage() (string, error) {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("  > docker ", "white", "", "bold") + "\t" +
		cmd.Colorfy(p.String(), "cyan", "", "")))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String()), nil
}

func (p *dockerProvisioner) GitDeploy(box *provision.Box, w io.Writer) (string, error) {
	imageId, err := p.gitDeploy(box.Repo, box.ImageVersion, w)
	if err != nil {
		return "", err
	}
	return p.deployPipeline(box, imageId, w)
}

func (p *dockerProvisioner) gitDeploy(re *repository.Repo, version string, w io.Writer) (string, error) {
	return p.getBuildImage(re, version), nil
}

func (p *dockerProvisioner) ImageDeploy(box *provision.Box, imageId string, w io.Writer) (string, error) {
	isValid, err := isValidBoxImage(box.GetFullName(), imageId)
	if err != nil {
		return "", err
	}
	if !isValid {
		return "", fmt.Errorf("invalid image for box %s: %s", box.GetFullName(), imageId)
	}
	return p.deployPipeline(box, imageId, w)
}

func (p *dockerProvisioner) deployPipeline(box *provision.Box, imageId string, w io.Writer) (string, error) {
	fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("--- deploy box (%s, image:%s)", box.GetFullName(), imageId)))
	actions := []*action.Action{
		&updateStatusInScylla,
		&createContainer,
		&updateContainerIdInScylla,
		&MileStoneUpdate,
		&updateStatusInScylla,
		&startContainer,
		&MileStoneUpdate,
		&updateStatusInScylla,
		&setNetworkInfo,
		&updateStatusInScylla,
		&followLogsAndCommit,
		&MileStoneUpdate,
		&updateStatusInScylla,
	}

	pipeline := action.NewPipeline(actions...)

	args := runContainerActionsArgs{
		box:             box,
		imageId:         imageId,
		writer:          w,
		isDeploy:        true,
		buildingImage:   imageId,
		containerState:  constants.StateInitializing,
		containerStatus: constants.StatusContainerLaunching,
		provisioner:     p,
	}
	err := pipeline.Execute(args)
	if err != nil {

		fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("deploy pipeline for box (%s) --> %s", box.GetFullName(), err)))
		return "", err
	}
	return imageId, nil
}

func (p *dockerProvisioner) Destroy(box *provision.Box, w io.Writer) error {

	fmt.Fprintf(w, lb.W(lb.DESTORYING, lb.INFO, fmt.Sprintf("\n--- destroying box (%s) ----", box.GetFullName())))
	containers, err := p.listContainersByBox(box)
	if err != nil {

		fmt.Fprintf(w, lb.W(lb.DESTORYING, lb.ERROR, fmt.Sprintf("Failed to list box containers (%s) --> %s", box.GetFullName(), err)))
		return err
	}
	args := changeUnitsPipelineArgs{
		box:         box,
		toRemove:    containers,
		writer:      ioutil.Discard,
		provisioner: p,
		boxDestroy:  true,
	}
	pipeline := action.NewPipeline(
		&destroyOldContainers,
		&removeOldRoutes,
	)
	err = pipeline.Execute(args)
	if err != nil {
		return err
	}
	return nil
}

func (p *dockerProvisioner) Start(box *provision.Box, process string, w io.Writer) error {
	containers, err := p.listContainersByBox(box)
	if err != nil {

		fmt.Fprintf(w, lb.W(lb.STARTING, lb.ERROR, fmt.Sprintf("Failed to list box containers (%s) --> %s", box.GetFullName(), err)))
	}
	return runInContainers(containers, func(c *container.Container, _ chan *container.Container) error {
		err := c.Start(&container.StartArgs{
			Provisioner: p,
			Box:         box,
		})
		if err != nil {
			return err
		}
		c.SetStatus(constants.StatusContainerStarting)
		if info, err := c.NetworkInfo(p); err == nil {
			p.fixContainer(c, info)
		}
		return nil
	}, nil, true)
}

func (p *dockerProvisioner) Stop(box *provision.Box, process string, w io.Writer) error {
	containers, err := p.listContainersByBox(box)
	if err != nil {
		fmt.Fprintf(w, lb.W(lb.STOPPING, lb.ERROR, fmt.Sprintf("Failed to list box containers (%s) --> %s", box.GetFullName(), err)))
	}
	return runInContainers(containers, func(c *container.Container, _ chan *container.Container) error {
		err := c.Stop(p)
		if err != nil {
			log.Errorf("Failed to stop %q: %s", box.GetFullName(), err)
		}
		return err
	}, nil, true)
}

func (p *dockerProvisioner) Restart(box *provision.Box, process string, w io.Writer) error {
	return nil
}

func (*dockerProvisioner) Addr(box *provision.Box) (string, error) {
	r, err := getRouterForBox(box)
	if err != nil {
		log.Errorf("Failed to get router: %s", err)
		return "", err
	}
	addr, err := r.Addr(box.GetFullName())
	if err != nil {
		log.Errorf("Failed to obtain box %s address: %s", box.GetFullName(), err)
		return "", err
	}
	return addr, nil
}

func (p *dockerProvisioner) SetBoxStatus(box *provision.Box, w io.Writer, status utils.Status) error {

	fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("---- status %s box %s ----", box.GetFullName(), status.String())))
	actions := []*action.Action{
		&updateStatusInScylla,
	}
	pipeline := action.NewPipeline(actions...)

	args := runContainerActionsArgs{
		box:             box,
		writer:          w,
		containerStatus: status,
		provisioner:     p,
	}

	err := pipeline.Execute(args)
	if err != nil {
		log.Errorf("error on execute status pipeline for box %s - %s", box.GetFullName(), err)
		return err
	}
	return nil
}

func (p *dockerProvisioner) SetCName(box *provision.Box, cname string) error {
	r, err := getRouterForBox(box)
	if err != nil {
		return err
	}
	return r.SetCName(cname, box.PublicIp)
}

func (p *dockerProvisioner) UnsetCName(box *provision.Box, cname string) error {
	r, err := getRouterForBox(box)
	if err != nil {
		return err
	}
	return r.UnsetCName(cname, box.PublicIp)
}

// PlatformAdd build and push a new docker platform to register
func (p *dockerProvisioner) PlatformAdd(name string, args map[string]string, w io.Writer) error {
	if args["dockerfile"] == "" {
		return errors.New("Dockerfile is required.")
	}
	if _, err := url.ParseRequestURI(args["dockerfile"]); err != nil {
		return errors.New("dockerfile parameter should be an url.")
	}
	imageName := platformImageName(name)
	cluster := p.Cluster()
	buildOptions := docker.BuildImageOptions{
		Name:           imageName,
		NoCache:        true,
		RmTmpContainer: true,
		Remote:         args["dockerfile"],
		InputStream:    nil,
		OutputStream:   w,
	}
	err := cluster.BuildImage(buildOptions)
	if err != nil {
		return err
	}
	parts := strings.Split(imageName, ":")
	var tag string
	if len(parts) > 2 {
		imageName = strings.Join(parts[:len(parts)-1], ":")
		tag = parts[len(parts)-1]
	} else if len(parts) > 1 {
		imageName = parts[0]
		tag = parts[1]
	} else {
		imageName = parts[0]
		tag = "latest"
	}
	return p.PushImage(imageName, tag)
}

func (p *dockerProvisioner) PlatformUpdate(name string, args map[string]string, w io.Writer) error {
	return p.PlatformAdd(name, args, w)
}

func (p *dockerProvisioner) PlatformRemove(name string) error {
	err := p.Cluster().RemoveImage(platformImageName(name))
	if err != nil && err == docker.ErrNoSuchImage {
		log.Errorf("error on remove image %s from docker.", name)
		return nil
	}
	return err
}

func (p *dockerProvisioner) Shell(opts provision.ShellOptions) error {
	var (
		c   *container.Container
		err error
	)
	c, err = p.GetContainerByBox(opts.Box)

	if err != nil {
		return err
	}
	return c.Shell(p, opts.Conn, opts.Conn, opts.Conn, container.Pty{Width: opts.Width, Height: opts.Height, Term: opts.Term})
}

func (p *dockerProvisioner) ExecuteCommandOnce(stdout, stderr io.Writer, box *provision.Box, cmd string, args ...string) error {
	container, err := p.GetContainerByBox(box)
	if err != nil {
		return err
	}
	return container.Exec(p, stdout, stderr, cmd, args...)
}

func (p *dockerProvisioner) MetricEnvs(start, end int64, point string, w io.Writer) ([]interface{}, error) {
	fmt.Fprintf(w, lb.W(lb.BILLING, lb.INFO, fmt.Sprintf("\n--- metrics collect for node (%s) ----", point)))
	res, err := p.Cluster().Showback(start, end, point)
	if err != nil {
		fmt.Fprintf(w, lb.W(lb.BILLING, lb.ERROR, fmt.Sprintf("--- pull metrics for the duration error(%d, %d)-->%s", start, end)))
		return nil, err
	}

	fmt.Fprintf(w, lb.W(lb.BILLING, lb.INFO, fmt.Sprintf("--- pull metrics for the duration (%d, %d)OK", start, end)))
	return res, nil
}

func (p *dockerProvisioner) SaveImage(box *provision.Box, w io.Writer) error {
	return nil
}

func (p *dockerProvisioner) DeleteImage(box *provision.Box, w io.Writer) error {
	return nil
}

func (p *dockerProvisioner) AttachDisk(box *provision.Box, w io.Writer) error {
	return nil
}

func (p *dockerProvisioner) DetachDisk(box *provision.Box, w io.Writer) error {
	return nil
}

func (p *dockerProvisioner) TriggerBills(account_id, cat_id, name string) error {
	cont := &container.Container{
		Name:      name,
		CartonId:  cat_id,
		AccountId: account_id,
	}
	err := cont.Deduct()
	if err != nil {
		return err
	}
	return nil
}
