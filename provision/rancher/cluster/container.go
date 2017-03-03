package cluster

import (
	//	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/go-rancher/v2"
	//	"github.com/megamsys/libgo/cmd"
	constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/vertice/carton"
	//	"github.com/megamsys/vertice/metrix"
	"net"
	"net/url"
	//"sync"
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
/*func (c *Cluster) CreateContainer(opts docker.CreateContainerOptions) (string, *docker.Container, error) {
	return c.CreateContainerSchedulerOpts(opts)
}
*/
// Similar to CreateContainer but allows arbritary options to be passed to
// the scheduler.
func (c *Cluster) CreateContainerSchedulerOpts(opts client.Container) (string, *client.Container, error) {

	var (
		addr      string
		container *client.Container
		err       error
	)

	maxTries := 5
	for ; maxTries > 0; maxTries-- {
		node, err := c.getNodeClient(c.Region)
		if err != nil {
			return "", nil, err
		}
		addr = node.addr
		container, err = node.RancherClient.Container.Create(&opts)
		if err == nil {
			c.handleNodeSuccess(addr)
			break
		} else {
			log.Errorf("Error trying to create container in node %q: %s. Trying again in another node...", addr, err.Error())
			shouldIncrementFailures := false
			if nodeErr, ok := err.(RancherNodeError); ok {
				baseErr := nodeErr.BaseError()
				if urlErr, ok := baseErr.(*url.Error); ok {
					baseErr = urlErr.Err
				}
				_, isNetErr := baseErr.(*net.OpError)
				if isNetErr || nodeErr.cmd == "createContainer" {
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
	//err = c.storage().StoreContainer(container.ID, addr)
	//err = c.storage().StoreContainerByName(container.ID, container.Name)
	return addr, container, err
}

func (c *Cluster) GetContainerById(id string) (*client.Container, error) {
	node, err := c.getNodeClient(c.Region)
	if err != nil {
		return nil, err
	}
	return node.RancherClient.Container.ById(id)
}

func (c *Cluster) getContainerNode(hostId string) (*client.Host, error) {
	node, err := c.getNodeClient(c.Region)
	if err != nil {
		return nil, err
	}
	return node.RancherClient.Host.ById(hostId)
}

func (c *Cluster) getNodeClient(region string) (node, error) {
	var n node
	var addr, aid, access, secret string
	nodes, err := c.Nodes()
	if err != nil {
		return n, err
	}
	for _, v := range nodes {
		if v.Metadata[RANCHER_ZONE] == region {
			addr = v.Address
			aid = v.Metadata[ADMIN_ID]
			access = v.Metadata[ACCESSKEY]
			secret = v.Metadata[SECRETKEY]
		}
	}
	if addr == "" {
		return n, errors.New("selected region unavailable [" + c.Region + "]")
	}
	cliaddr := client.ClientOpts{Url: addr, AccountId: aid, AccessKey: access, SecretKey: secret}
	return c.getNodeByAddr(cliaddr)
}

func (c *Cluster) SetNetworkinNode(hostId, IpAddress, cartonId, email string) error {
	if hostId == "" {
		return errors.New("empty host Id")
	}
	host, err := c.getContainerNode(hostId)
	if err != nil {
		return err
	}
	err = c.IpNode(IpAddress, host.AgentIpAddress, cartonId, email)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cluster) IpNode(contIp, nodeIp, CartonId, email string) error {
	var ips = make(map[string][]string)
	ips[c.getIps()] = []string{contIp}
	ips[carton.HOSTIP] = []string{nodeIp}
	if asm, err := carton.NewAssembly(CartonId, email, ""); err != nil {
		return err
	} else if err = asm.NukeAndSetOutputs(ips); err != nil {
		return err
	}
	return nil
}

func (c *Cluster) getIps() string {
	for k, v := range c.VNets {
		if v == "true" {
			return k
		}
	}
	return constants.PRIVATEIPV4
}

// RemoveContainer removes a container from the cluster.
func (c *Cluster) RemoveContainer(opts *client.Container) error {
	return c.removeFromStorage(opts)
}

func (c *Cluster) removeFromStorage(opts *client.Container) error {
	node, err := c.getNodeClient(c.Region)
	if err != nil {
		return err
	}
	err = node.RancherClient.Container.Delete(opts)
	if err != nil {
		return wrapError(node, err)
	}
	return nil
}

func (c *Cluster) StartContainer(id string) error {
	node, err := c.getNodeClient(c.Region)
	if err != nil {
		return err
	}
	cont, err := node.RancherClient.Container.ById(id)
	if err != nil {
		return err
	}
	_, err = node.RancherClient.Container.ActionStart(cont)
	return wrapError(node, err)
}

func (c *Cluster) StopContainer(id string) error {
	node, err := c.getNodeClient(c.Region)
	if err != nil {
		return err
	}
	cont, err := node.RancherClient.Container.ById(id)
	if err != nil {
		return err
	}
	insStop, err := node.RancherClient.InstanceStop.ById(id)
	if err != nil {
		return err
	}

	_, err = node.RancherClient.Container.ActionStop(cont, insStop)
	return wrapError(node, err)
}
