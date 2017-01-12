package cluster

import (
//	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/client"
//	"github.com/megamsys/libgo/cmd"
//	constants "github.com/megamsys/libgo/utils"
	//"github.com/megamsys/vertice/carton"
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
		nodes, err := c.Nodes()
		for _, v := range nodes {
			if v.Metadata[RANCHER_ZONE] == c.Region {
				addr = v.Address
			}
		}
		if addr == "" {
			return addr, nil, errors.New("CreateContainer needs a non empty node addr")
		}

     cliaddr := client.ClientOpts{ Url: addr ,}
		node, err := c.getNodeByAddr(cliaddr)
		if err != nil {
			return addr,nil, err
		}

		container, err= node.RancherClient.Container.Create(&opts)
	  fmt.Println(container)

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
//	err = c.storage().StoreContainerByName(container.ID, container.Name)
	return addr, container, err
}
