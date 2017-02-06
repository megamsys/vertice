package cluster

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/megamsys/libgo/cmd"
	constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/vertice/carton"
	"github.com/megamsys/vertice/metrix"
	"net"
	"net/url"
	"sync"
	//	"time"
)

type Container struct {
	Id   string
	Host string
}

// CreateContainer creates a container in the specified node. If no node is
// specified, it will create the container in a node selected by the scheduler.
//
// It returns the container, or an error, in case of failures.
func (c *Cluster) CreateContainer(opts docker.CreateContainerOptions) (string, *docker.Container, error) {
	return c.CreateContainerSchedulerOpts(opts)
}

// Similar to CreateContainer but allows arbritary options to be passed to
// the scheduler.
func (c *Cluster) CreateContainerSchedulerOpts(opts docker.CreateContainerOptions) (string, *docker.Container, error) {
	var (
		addr      string
		container *docker.Container
		err       error
	)

	maxTries := 5
	for ; maxTries > 0; maxTries-- {
		nodes, err := c.Nodes()
		for _, v := range nodes {
			if v.Metadata[DOCKER_ZONE] == c.Region {
				addr = v.Address
			}
		}
		if addr == "" {
			return addr, nil, errors.New("CreateContainer needs a non empty node addr")
		}
		container, err = c.createContainerInNode(opts, addr)
		if err == nil {
			c.handleNodeSuccess(addr)
			break
		} else {
			log.Errorf("Error trying to create container in node %q: %s. Trying again in another node...", addr, err.Error())
			shouldIncrementFailures := false
			if nodeErr, ok := err.(DockerNodeError); ok {
				baseErr := nodeErr.BaseError()
				if urlErr, ok := baseErr.(*url.Error); ok {
					baseErr = urlErr.Err
				}
				_, isNetErr := baseErr.(*net.OpError)
				if isNetErr || baseErr == docker.ErrConnectionRefused || nodeErr.cmd == "createContainer" {
					shouldIncrementFailures = true
				}
			}
			c.handleNodeError(addr, err, shouldIncrementFailures)
			return addr, nil, err
		}
	}
	if err != nil {
		return addr, nil, fmt.Errorf("CreateContainer: maximum number of tries exceeded, last error: %s", err.Error())
	}
	err = c.storage().StoreContainer(container.ID, addr)
	err = c.storage().StoreContainerByName(container.ID, container.Name)
	return addr, container, err
}

func (c *Cluster) createContainerInNode(opts docker.CreateContainerOptions, nodeAddress string) (*docker.Container, error) {
	registryServer, _ := parseImageRegistry(opts.Config.Image)
	if registryServer != "" {
		err := c.PullImage(docker.PullImageOptions{
			Repository: opts.Config.Image,
		}, docker.AuthConfiguration{}, nodeAddress)
		if err != nil {
			return nil, err
		}
	}
	node, err := c.getNodeByAddr(nodeAddress)
	if err != nil {
		return nil, err
	}
	cont, err := node.CreateContainer(opts)
	return cont, wrapErrorWithCmd(node, err, "createContainer")
}

// InspectContainer returns information about a container by its ID, getting
// the information from the right node.
func (c *Cluster) InspectContainer(id string) (*docker.Container, error) {
	node, err := c.getNodeForContainer(id)
	if err != nil {
		return nil, err
	}
	cont, err := node.InspectContainer(id)
	return cont, wrapError(node, err)
}

// KillContainer kills a container, returning an error in case of failure.
func (c *Cluster) KillContainer(opts docker.KillContainerOptions) error {
	node, err := c.getNodeForContainer(opts.ID)
	if err != nil {
		return err
	}
	return wrapError(node, node.KillContainer(opts))
}

// ListContainers returns a slice of all containers in the cluster matching the
// given criteria.
func (c *Cluster) ListContainers(opts docker.ListContainersOptions) ([]docker.APIContainers, error) {
	var addr string
	nodes, err := c.Nodes()
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	result := make(chan []docker.APIContainers, len(nodes))
	errs := make(chan error, len(nodes))
	for _, n := range nodes {
		if n.Metadata[DOCKER_ZONE] == c.Region {
			addr = n.Address
		}
		wg.Add(1)
		client, _ := c.getNodeByAddr(addr)
		go func(n node) {
			defer wg.Done()
			if containers, err := n.ListContainers(opts); err != nil {
				errs <- wrapError(n, err)
			} else {
				result <- containers
			}
		}(client)
	}
	wg.Wait()
	var group []docker.APIContainers
	for {
		select {
		case containers := <-result:
			group = append(group, containers...)
		case err = <-errs:
		default:
			return group, err
		}
	}
}

// RemoveContainer removes a container from the cluster.
func (c *Cluster) RemoveContainer(opts docker.RemoveContainerOptions) error {
	return c.removeFromStorage(opts)
}

func (c *Cluster) removeFromStorage(opts docker.RemoveContainerOptions) error {
	node, err := c.getNodeForContainer(opts.ID)
	if err != nil {
		return err
	}
	err = node.RemoveContainer(opts)
	if err != nil {
		_, isNoSuchContainer := err.(*docker.NoSuchContainer)
		if !isNoSuchContainer {
			return wrapError(node, err)
		}
	}
	return c.storage().RemoveContainer(opts.ID)
}

func (c *Cluster) StartContainer(id string, hostConfig *docker.HostConfig) error {

	var n node
	n, err := c.getNodeForContainer(id)
	if err != nil {
		n, err = c.getNodeByRegion(c.Region)
		if err != nil {
		  return err
	  }
	}
	return wrapError(n, n.StartContainer(id, hostConfig))
}

func (c *Cluster) PreStopAction(name string) (string, error) {
	id, err := c.storage().RetrieveContainerByName(name)
	if err != nil {
		return "", err
	}
	return id, err
}

// StopContainer stops a container, killing it after the given timeout, if it
// fails to stop nicely.
func (c *Cluster) StopContainer(id string, timeout uint) error {
	var n node
	n, err := c.getNodeForContainer(id)
	if err != nil {
		n, err = c.getNodeByRegion(c.Region)
		if err != nil {
		  return err
	  }
	}
	return wrapError(n, n.StopContainer(id, timeout))
}

// RestartContainer restarts a container, killing it after the given timeout,
// if it fails to stop nicely.
func (c *Cluster) RestartContainer(id string, timeout uint) error {
	node, err := c.getNodeForContainer(id)
	if err != nil {
		return err
	}
	return wrapError(node, node.RestartContainer(id, timeout))
}

// PauseContainer changes the container to the paused state.
func (c *Cluster) PauseContainer(id string) error {
	node, err := c.getNodeForContainer(id)
	if err != nil {
		return err
	}
	return wrapError(node, node.PauseContainer(id))
}

// UnpauseContainer removes the container from the paused state.
func (c *Cluster) UnpauseContainer(id string) error {
	node, err := c.getNodeForContainer(id)
	if err != nil {
		return err
	}
	return wrapError(node, node.UnpauseContainer(id))
}

// WaitContainer blocks until the given container stops, returning the exit
// code of the container command.
func (c *Cluster) WaitContainer(id string) (int, error) {
	node, err := c.getNodeForContainer(id)
	if err != nil {
		return -1, err
	}
	code, err := node.WaitContainer(id)
	return code, wrapError(node, err)
}

// AttachToContainer attaches to a container, using the given options.
func (c *Cluster) AttachToContainer(opts docker.AttachToContainerOptions) error {
	node, err := c.getNodeForContainer(opts.Container)
	if err != nil {
		return err
	}
	return wrapError(node, node.AttachToContainer(opts))
}

// Logs retrieves the logs of the specified container.
func (c *Cluster) Logs(opts docker.LogsOptions) error {
	node, err := c.getNodeForContainer(opts.Container)
	if err != nil {
		return err
	}
	return wrapError(node, node.Logs(opts))
}

// CommitContainer commits a container and returns the image id.
func (c *Cluster) CommitContainer(opts docker.CommitContainerOptions) (*docker.Image, error) {
	node, err := c.getNodeForContainer(opts.Container)
	if err != nil {
		return nil, err
	}
	image, err := node.CommitContainer(opts)
	if err != nil {
		return nil, wrapError(node, err)
	}
	key := imageKey(opts.Repository, opts.Tag)
	if key != "" {
		err = c.storage().StoreImage(key, image.ID, node.addr)
		if err != nil {
			return nil, err
		}
	}
	return image, nil
}

// ExportContainer exports a container as a tar and writes
// the result in out.
func (c *Cluster) ExportContainer(opts docker.ExportContainerOptions) error {
	node, err := c.getNodeForContainer(opts.ID)
	if err != nil {
		return err
	}
	return wrapError(node, node.ExportContainer(opts))
}

// TopContainer returns information about running processes inside a container
// by its ID, getting the information from the right node.
func (c *Cluster) TopContainer(id string, psArgs string) (docker.TopResult, error) {
	node, err := c.getNodeForContainer(id)
	if err != nil {
		return docker.TopResult{}, err
	}
	result, err := node.TopContainer(id, psArgs)
	return result, wrapError(node, err)
}

func (c *Cluster) getNodeForContainer(container string) (node, error) {
	return c.getNode(func(s Storage) (string, error) {
		return s.RetrieveContainer(container)
	})
}

func (c *Cluster) SetNetworkinNode(containerId, cartonId, email string) error {
	port := c.GulpPort()
	container := c.getContainerObject(containerId)
	err := c.Ips(container.NetworkSettings.IPAddress, cartonId,email)
	if err != nil {
		return err
	}
	client := DockerClient{ContainerId: containerId, CartonId: cartonId, AccountId: email}
	err = client.NetworkRequest(container.Node.IP, port)

	if err != nil {
		return err
	}
	return nil
}

func (c *Cluster) Ips(ip, CartonId,email string) error {
	var ips = make(map[string][]string)
	pubipv4s := []string{ip}
	ips[c.getIps()] = pubipv4s
	if asm, err := carton.NewAssembly(CartonId,email, ""); err != nil {
		return err
	} else if err = asm.NukeAndSetOutputs(ips); err != nil {
		return err
	}
	return nil
}

func (c *Cluster) getIps() string {
	for k, v := range c.VNets {
		if v == "true" {
			switch k {
			case constants.IPV4PUB:
				return carton.PUBLICIPV4
			case constants.IPV6PUB:
				return carton.PUBLICIPV6
			case constants.IPV4PRI:
				return carton.PRIVATEIPV4
			case constants.IPV6PRI:
        return carton.PRIVATEIPV6
			}
		}
	}
	return ""
}

func (c *Cluster) SetLogs(cs chan []byte, opts docker.LogsOptions, closechan chan bool) error {
	node, err := c.getNodeForContainer(opts.Container)
	node.Logs(opts)
	closechan <- true
	if err != nil {
		return err
	}
	return nil
}

func (c *Cluster) getContainerObject(containerId string) *docker.Container {
	inspect, _ := c.InspectContainer(containerId) //gets the swarmNode

	container := &docker.Container{}
	insp, _ := json.Marshal(inspect)
	json.Unmarshal([]byte(string(insp)), container)

	return container

}

func (c *Cluster) CreateExec(opts docker.CreateExecOptions) (*docker.Exec, error) {
	node, err := c.getNodeForContainer(opts.Container)
	if err != nil {
		node, err = c.getNodeByRegion(c.Region)
		if err != nil {
		  return nil, err
	  }
	}
	exec, err := node.CreateExec(opts)
	return exec, wrapError(node, err)
}

func (c *Cluster) getNodeByRegion(region string) (node, error) {
 var addr	string
 var n node
	nodes, err := c.Nodes()
	if err != nil {
		return n, err
	}
	for _, v := range nodes {
		if v.Metadata[DOCKER_ZONE] == region {
			addr = v.Address
		}
	}
	if addr == "" {
		 return n, errors.New("CreateContainer needs a non empty node addr")
	}
 return c.getNodeByAddr(addr)
}

func (c *Cluster) StartExec(execId, containerId string, opts docker.StartExecOptions) error {
	node, err := c.getNodeForContainer(containerId)
	if err != nil {
		node, err = c.getNodeByRegion(c.Region)
		if err != nil {
		  return err
	  }
	}
	return wrapError(node, node.StartExec(execId, opts))
}

func (c *Cluster) ResizeExecTTY(execId, containerId string, height, width int) error {
	node, err := c.getNodeForContainer(containerId)
	if err != nil {
		node, err = c.getNodeByRegion(c.Region)
		if err != nil {
		  return err
	  }
	}
	return wrapError(node, node.ResizeExecTTY(execId, height, width))
}

func (c *Cluster) InspectExec(execId, containerId string) (*docker.ExecInspect, error) {
	node, err := c.getNodeForContainer(containerId)
	if err != nil {
		node, err = c.getNodeByRegion(c.Region)
		if err != nil {
		  return nil, err
	  }
	}
	execInspect, err := node.InspectExec(execId)
	if err != nil {
		return nil, wrapError(node, err)
	}
	return execInspect, nil
}
func (c *Cluster) GulpPort() string {
	var gulpPort string
	nodes, _ := c.Nodes()
	for _, v := range nodes {
		if v.Metadata[DOCKER_ZONE] == c.Region {
			gulpPort = v.Metadata[DOCKER_GULP]
		}
	}
	return gulpPort
}

// Showback returns the metrics of the swarm containers stats

func (c *Cluster) Showback(start int64, end int64, point string) ([]interface{}, error) {
	log.Debugf("showback (%d, %d)", start, end)
	var (
		result *docker.Container
		v  docker.APIContainers
		resultStats []interface{}
	)
	node, err := c.getNodeByAddr(point)
	if err != nil {
		return nil, fmt.Errorf("%s", cmd.Colorfy("Unavailable nodes (hint: start or beat it).\n", "red", "", ""))
	}
	opts := docker.ListContainersOptions{
		All: true,
		//	Filters: map[string][]string{"status": {"running","paused","stopped"}},
	}
	ps, err := node.ListContainers(opts)
	if err != nil {
		return nil, err
	}
	for _, v = range ps {
		id := v.ID
		result, _ = node.InspectContainer(id)
		res := &metrix.Stats{
			ContainerId:  result.ID,
			Image: result.Image,
			AllocatedMemory: result.HostConfig.Memory,
			AllocatedCpu: result.HostConfig.CPUShares,
			AccountId:    v.Labels[constants.ACCOUNT_ID],
			AssemblyId:   v.Labels[constants.ASSEMBLY_ID],
			AssembliesId: v.Labels[constants.ASSEMBLIES_ID],
			AssemblyName: v.Labels[constants.ASSEMBLY_NAME],
			CPUUnitCost: v.Labels[carton.CONTAINER_CPU_COST],
			MemoryUnitCost: v.Labels[carton.CONTAINER_MEMORY_COST],
			QuotaId: v.Labels[constants.QUOTA_ID],
			//AuditPeriod:  stats.Read,
			Status:       v.State,
		}
		resultStats = append(resultStats,res)
	}
		return resultStats, nil
		}
