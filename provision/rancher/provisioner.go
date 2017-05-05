package rancher

import (
	"bytes"
	//	"errors"
	"fmt"
	"io"
	"io/ioutil"
	//	"net/url"
	"strings"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	//"github.com/megamsys/go-rancher/v2"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/libgo/utils"
	constants "github.com/megamsys/libgo/utils"
	lb "github.com/megamsys/vertice/logbox"
	"github.com/megamsys/vertice/provision"
	"github.com/megamsys/vertice/provision/rancher/cluster"
	"github.com/megamsys/vertice/provision/rancher/container"
	"github.com/megamsys/vertice/repository"
	"github.com/megamsys/vertice/router"
	_ "github.com/megamsys/vertice/router/route53"
	"github.com/megamsys/vertice/toml"
)

var mainRancherProvisioner *rancherProvisioner

func init() {
	mainRancherProvisioner = &rancherProvisioner{}
	provision.Register("rancher", mainRancherProvisioner)
}

type rancherProvisioner struct {
	cluster        *cluster.Cluster
	collectionName string
	storage        cluster.Storage
}
type Rancher struct {
	Enabled bool     `json:"enabled" toml:"enabled"`
	Regions []Region `json:"region" toml:"region"`
}

type Region struct {
	RancherZone     string        `json:"rancher_zone" toml:"rancher_zone"`
	RancherEndPoint string        `json:"rancher" toml:"rancher"`
	Registry        string        `json:"registry" toml:"registry"`
	CPUPeriod       toml.Duration `json:"cpu_period" toml:"cpu_period"`
	CPUQuota        toml.Duration `json:"cpu_quota" toml:"cpu_quota"`
	AdminId         string        `json:"admin_id" toml:"admin_id"`
	AdminAccess     string        `json:"access_key" toml:"access_key"`
	AdminSecret     string        `json:"secret_key" toml:"secret_key"`
}

func (p *rancherProvisioner) Cluster() *cluster.Cluster {
	if p.cluster == nil {
		panic("✗ rancher cluster")
	}
	return p.cluster
}

func (p *rancherProvisioner) String() string {
	if p.cluster == nil {
		return "✗ rancher cluster"
	}
	return "ready"
}

func (p *rancherProvisioner) Initialize(m interface{}) error {
	return p.initRancherCluster(m)
}

func (p *rancherProvisioner) initRancherCluster(i interface{}) error {
	var err error
	if p.storage == nil {
		p.storage, err = buildClusterStorage()
		if err != nil {
			return err
		}
	}
	if w, ok := i.(Rancher); ok {
		var nodes []cluster.Node
		for i := 0; i < len(w.Regions); i++ {
			m := w.Regions[i].toMap()
			n := cluster.Node{
				Address:  m[cluster.RANCHER_SERVER], //rancher endpoint
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
	m[cluster.RANCHER_ZONE] = c.RancherZone
	m[cluster.RANCHER_SERVER] = c.RancherEndPoint
	m[cluster.ADMIN_ID] = c.AdminId
	m[cluster.ACCESSKEY] = c.AdminAccess
	m[cluster.SECRETKEY] = c.AdminSecret
	m[cluster.RANCHER_REGISTRY] = c.Registry
	m[cluster.RANCHER_CPUPERIOD] = c.CPUPeriod.String()
	m[cluster.RANCHER_CPUQUOTA] = c.CPUQuota.String()
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

func (p *rancherProvisioner) StartupMessage() (string, error) {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("  > rancher ", "white", "", "bold") + "\t" +
		cmd.Colorfy(p.String(), "cyan", "", "")))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String()), nil
}

func (p *rancherProvisioner) GitDeploy(box *provision.Box, w io.Writer) (string, error) {
	imageId, err := p.gitDeploy(box.Repo, box.ImageVersion, w)
	if err != nil {
		return "", err
	}
	return p.deployPipeline(box, imageId, w)
}

func (p *rancherProvisioner) gitDeploy(re *repository.Repo, version string, w io.Writer) (string, error) {
	return p.getBuildImage(re, version), nil
	//return "",nil
}

func (p *rancherProvisioner) ImageDeploy(box *provision.Box, imageId string, w io.Writer) (string, error) {
	isValid, err := isValidBoxImage(box.GetFullName(), imageId)
	if err != nil {
		return "", err
	}
	if !isValid {
		return "", fmt.Errorf("invalid image for box %s: %s", box.GetFullName(), imageId)
	}
	return p.deployPipeline(box, imageId, w)
}

func (p *rancherProvisioner) BackupDeploy(box *provision.Box, imageId string, w io.Writer) (string, error) {
	return "", nil
}

func (p *rancherProvisioner) deployPipeline(box *provision.Box, imageId string, w io.Writer) (string, error) {
	fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("--- deploy box (%s, image:%s)", box.GetFullName(), imageId)))
	actions := []*action.Action{
		&updateStatusInScylla,
		&createContainer,
		&updateContainerIdInScylla,
		&MileStoneUpdate,
		&updateStatusInScylla,
		&waitToContainerUp,
		&MileStoneUpdate,
		&updateStatusInScylla,
		&setNetworkInfo,
		&updateStatusInScylla,
		//	&followLogsAndCommit,
		//	&MileStoneUpdate,
		//	&updateStatusInScylla,
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

func (p *rancherProvisioner) Destroy(box *provision.Box, w io.Writer) error {

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
	//&removeOldRoutes,
	)
	err = pipeline.Execute(args)
	if err != nil {
		return err
	}
	return nil
}

func (p *rancherProvisioner) Start(box *provision.Box, process string, w io.Writer) error {
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
		return nil
	}, nil, true)
}

func (p *rancherProvisioner) Stop(box *provision.Box, process string, w io.Writer) error {
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

func (p *rancherProvisioner) Restart(box *provision.Box, process string, w io.Writer) error {
	return nil
}

func (*rancherProvisioner) Addr(box *provision.Box) (string, error) {

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

func (p *rancherProvisioner) SetBoxStatus(box *provision.Box, w io.Writer, status utils.Status) error {

	fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("---- status %s box %s ----", box.GetFullName(), status.String())))
	actions := []*action.Action{
	//&updateStatusInScylla,
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

func (p *rancherProvisioner) SetCName(box *provision.Box, cname string) error {
	r, err := getRouterForBox(box)
	if err != nil {
		return err
	}
	return r.SetCName(cname, box.PublicIp)
}

func (p *rancherProvisioner) UnsetCName(box *provision.Box, cname string) error {
	r, err := getRouterForBox(box)
	if err != nil {
		return err
	}
	return r.UnsetCName(cname, box.PublicIp)
}

// PlatformAdd build and push a new docker platform to register
func (p *rancherProvisioner) PlatformAdd(name string, args map[string]string, w io.Writer) error {
	/*
		if args["dockerfile"] == "" {
			return errors.New("Rancherfile is required.")
		}
		if _, err := url.ParseRequestURI(args["rancherfile"]); err != nil {
			return errors.New("rancherfile parameter should be an url.")
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
		}*/
	//return p.PushImage("", tag)

	return nil
}

func (p *rancherProvisioner) PlatformUpdate(name string, args map[string]string, w io.Writer) error {
	return p.PlatformAdd(name, args, w)
}

func (p *rancherProvisioner) PlatformRemove(name string) error {
	/*
		err := p.Cluster().RemoveImage(platformImageName(name))
		if err != nil && err == docker.ErrNoSuchImage {
			log.Errorf("error on remove image %s from docker.", name)
			return nil
		}
		return err
	*/
	return nil
}

func (p *rancherProvisioner) Shell(opts provision.ShellOptions) error {
	var (
		err error
	)
	_, err = p.GetContainerByBox(opts.Box)

	if err != nil {
		return err
	}
	//return c.Shell(p, opts.Conn, opts.Conn, opts.Conn, container.Pty{Width: opts.Width, Height: opts.Height, Term: opts.Term})
	return nil
}

func (p *rancherProvisioner) ExecuteCommandOnce(stdout, stderr io.Writer, box *provision.Box, cmd string, args ...string) error {
	_, err := p.GetContainerByBox(box)
	if err != nil {
		return err
	}
	//return container.Exec(p, stdout, stderr, cmd, args...)
	return nil
}

func (p *rancherProvisioner) MetricEnvs(start, end int64, point string, w io.Writer) ([]interface{}, error) {
	/*
		fmt.Fprintf(w, lb.W(lb.BILLING, lb.INFO, fmt.Sprintf("\n--- metrics collect for node (%s) ----", point)))
		//res, err := p.Cluster().Showback(start, end, point)
		if err != nil {
			fmt.Fprintf(w, lb.W(lb.BILLING, lb.ERROR, fmt.Sprintf("--- pull metrics for the duration error(%d, %d)-->%s", start, end)))
			return nil, err
		}

		fmt.Fprintf(w, lb.W(lb.BILLING, lb.INFO, fmt.Sprintf("--- pull metrics for the duration (%d, %d)OK", start, end)))
		//return res, nil
	*/
	var b []interface{}
	//a :=[]string{}
	//b = a
	return b, nil
}

func (p *rancherProvisioner) SaveImage(box *provision.Box, w io.Writer) error {
	return nil
}

func (p *rancherProvisioner) DeleteImage(box *provision.Box, w io.Writer) error {
	return nil
}

func (p *rancherProvisioner) CreateSnapshot(box *provision.Box, w io.Writer) error {
	return nil
}

func (p *rancherProvisioner) DeleteSnapshot(box *provision.Box, w io.Writer) error {
	return nil
}

func (p *rancherProvisioner) RestoreSnapshot(box *provision.Box, w io.Writer) error {
	return nil
}

func (p *rancherProvisioner) AttachDisk(box *provision.Box, w io.Writer) error {
	return nil
}

func (p *rancherProvisioner) DetachDisk(box *provision.Box, w io.Writer) error {
	return nil
}

func (p *rancherProvisioner) Suspend(box *provision.Box, process string, w io.Writer) error {
	return nil
}

func (p *rancherProvisioner) TriggerBills(account_id, cat_id, name string) error {
	/*
		cont := &container.Container{
			Name:      name,
			CartonId:  cat_id,
			AccountId: account_id,
		}
		err := cont.Deduct()
		if err != nil {
			return err
		}
	*/
	return nil
}
