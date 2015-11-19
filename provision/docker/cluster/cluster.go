package cluster

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/fsouza/go-dockerclient"
)

var (
	errStorageMandatory = errors.New("Storage parameter is mandatory")
	errHealerInProgress = errors.New("Healer already running")

	timeout10Client  = clientWithTimeout(10*time.Second, 1*time.Hour)
	persistentClient = clientWithTimeout(10*time.Second, 0)
	timeout10Dialer  = &net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}
)

type node struct {
	addr string
	*docker.Client
}

// ContainerStorage provides methods to store and retrieve information about
// the relation between the node and the container. It can be easily
// represented as a key-value storage.
//
// The relevant information is: in which host the given container is running?
type ContainerStorage interface {
	StoreContainer(container, host string) error
	RetrieveContainer(container string) (host string, err error)
	RemoveContainer(container string) error
	RetrieveContainers() ([]Container, error)

	StoreContainerByName(container, host string) error
	RetrieveContainerByName(name string) (container string, err error)
}

// ImageStorage works like ContainerStorage, but stores information about
// images and hosts.
type ImageStorage interface {
	StoreImage(repo, id, host string) error
	RetrieveImage(repo string) (Image, error)
	RemoveImage(repo, id, host string) error
	RetrieveImages() ([]Image, error)
}

type NodeStorage interface {
	StoreNode(node Node) error
	RetrieveNodesByMetadata(metadata map[string]string) ([]Node, error)
	RetrieveNodes() ([]Node, error)
	RetrieveNode(address string) (Node, error)
	UpdateNode(node Node) error
	RemoveNode(address string) error
	LockNodeForHealing(address string, isFailure bool, timeout time.Duration) (bool, error)
	ExtendNodeLock(address string, timeout time.Duration) error
	UnlockNode(address string) error
}

type Storage interface {
	ContainerStorage
	ImageStorage
	NodeStorage
}

// Cluster is the basic type of the package. It manages internal nodes, and
// provide methods for interaction with those nodes, like CreateContainer,
// which creates a container in one node of the cluster.
type Cluster struct {
	Healer         Healer
	stor           Storage
	bridges        Bridges
	gulp           Gulp
	monitoringDone chan bool
}

type DockerNodeError struct {
	node node
	cmd  string
	err  error
}

func (n DockerNodeError) Error() string {
	if n.cmd == "" {
		return fmt.Sprintf("error in docker node %q: %s", n.node.addr, n.err.Error())
	}
	return fmt.Sprintf("error in docker node %q running command %q: %s", n.node.addr, n.cmd, n.err.Error())
}

func (n DockerNodeError) BaseError() error {
	return n.err
}

func wrapError(n node, err error) error {
	if err != nil {
		return DockerNodeError{node: n, err: err}
	}
	return nil
}

func wrapErrorWithCmd(n node, err error, cmd string) error {
	if err != nil {
		return DockerNodeError{node: n, err: err, cmd: cmd}
	}
	return nil
}

// New creates a new Cluster, initially composed by the given nodes.
//
// The scheduler parameter defines the scheduling strategy. It defaults
// to round robin if nil.
// The storage parameter is the storage the cluster instance will use.
func New(storage Storage, gulp Gulp, bridges []Bridge, nodes ...Node) (*Cluster, error) {
	var (
		c   Cluster
		err error
	)
	if storage == nil {
		return nil, errStorageMandatory
	}
	c.stor = storage
	c.bridges = bridges
	c.gulp = gulp
	c.Healer = DefaultHealer{}

	if len(nodes) > 0 {
		for _, n := range nodes {
			err = c.Register(n)
			if err != nil {
				return &c, err
			}
		}
	}
	return &c, err
}

// Register adds new nodes to the cluster.
func (c *Cluster) Register(node Node) error {
	if node.Address == "" {
		return errors.New("Invalid address")
	}
	return c.storage().StoreNode(node)
}

func (c *Cluster) UpdateNode(node Node) (Node, error) {
	unlock, err := c.lockWithTimeout(node.Address, false)
	if err != nil {
		return Node{}, err
	}
	defer unlock()
	dbNode, err := c.storage().RetrieveNode(node.Address)
	if err != nil {
		return Node{}, err
	}
	if node.CreationStatus != "" && node.CreationStatus != dbNode.CreationStatus {
		if dbNode.CreationStatus != NodeCreationStatusPending && dbNode.CreationStatus != "" {
			return Node{}, fmt.Errorf("cannot update node status when current status is %q", dbNode.CreationStatus)
		}
		dbNode.CreationStatus = node.CreationStatus
	}
	for k, v := range node.Metadata {
		if v == "" {
			delete(dbNode.Metadata, k)
		} else {
			dbNode.Metadata[k] = v
		}
	}
	return dbNode, c.storage().UpdateNode(dbNode)
}

// Unregister removes nodes from the cluster.
func (c *Cluster) Unregister(address string) error {
	return c.storage().RemoveNode(address)
}

func (c *Cluster) UnfilteredNodes() ([]Node, error) {
	return c.storage().RetrieveNodes()
}

func (c *Cluster) Nodes() ([]Node, error) {
	nodes, err := c.storage().RetrieveNodes()
	if err != nil {
		return nil, err
	}
	return NodeList(nodes).filterDisabled(), nil
}

func (c *Cluster) NodesForMetadata(metadata map[string]string) ([]Node, error) {
	nodes, err := c.storage().RetrieveNodesByMetadata(metadata)
	if err != nil {
		return nil, err
	}
	return NodeList(nodes).filterDisabled(), nil
}

func (c *Cluster) lockWithTimeout(addr string, isFailure bool) (func(), error) {
	lockTimeout := 3 * time.Minute
	locked, err := c.storage().LockNodeForHealing(addr, isFailure, lockTimeout)
	if err != nil {
		return nil, err
	}
	if !locked {
		return nil, errHealerInProgress
	}
	doneKeepAlive := make(chan bool)
	go func() {
		for {
			select {
			case <-doneKeepAlive:
				return
			case <-time.After(30 * time.Second):
			}
			c.storage().ExtendNodeLock(addr, lockTimeout)
		}
	}()
	return func() {
		doneKeepAlive <- true
		c.storage().UnlockNode(addr)
	}, nil
}

func (c *Cluster) handleNodeError(addr string, lastErr error, incrementFailures bool) error {
	unlock, err := c.lockWithTimeout(addr, true)
	if err != nil {
		return err
	}
	go func() {
		defer unlock()
		node, err := c.storage().RetrieveNode(addr)
		if err != nil {
			return
		}
		node.updateError(lastErr, incrementFailures)
		duration := c.Healer.HandleError(&node)
		if duration > 0 {
			node.updateDisabled(time.Now().Add(duration))
		}
		c.storage().UpdateNode(node)
	}()
	return nil
}

func (c *Cluster) handleNodeSuccess(addr string) error {
	unlock, err := c.lockWithTimeout(addr, false)
	if err != nil {
		return err
	}
	defer unlock()
	node, err := c.storage().RetrieveNode(addr)
	if err != nil {
		return err
	}
	node.updateSuccess()
	return c.storage().UpdateNode(node)
}

func (c *Cluster) storage() Storage {
	return c.stor
}

type nodeFunc func(node) (interface{}, error)

func (c *Cluster) runOnNodes(fn nodeFunc, errNotFound error, wait bool, nodeAddresses ...string) (interface{}, error) {
	if len(nodeAddresses) == 0 {
		nodes, err := c.Nodes()
		if err != nil {
			return nil, err
		}
		nodeAddresses = make([]string, len(nodes))
		for i, node := range nodes {
			nodeAddresses[i] = node.Address
		}
	}
	var wg sync.WaitGroup
	finish := make(chan int8, len(nodeAddresses))
	errChan := make(chan error, len(nodeAddresses))
	result := make(chan interface{}, len(nodeAddresses))
	for _, addr := range nodeAddresses {
		wg.Add(1)
		client, err := c.getNodeByAddr(addr)
		if err != nil {
			return nil, err
		}
		go func(n node) {
			defer wg.Done()
			value, err := fn(n)
			if err == nil {
				result <- value
			} else if e, ok := err.(*docker.Error); ok && e.Status == http.StatusNotFound {
				return
			} else if !reflect.DeepEqual(err, errNotFound) {
				errChan <- wrapError(n, err)
			}
		}(client)
	}
	if wait {
		wg.Wait()
		select {
		case value := <-result:
			return value, nil
		case err := <-errChan:
			return nil, err
		default:
			return nil, errNotFound
		}
	}
	go func() {
		wg.Wait()
		close(finish)
	}()
	select {
	case value := <-result:
		return value, nil
	case err := <-errChan:
		return nil, err
	case <-finish:
		select {
		case value := <-result:
			return value, nil
		default:
			return nil, errNotFound
		}
	}
}

func (c *Cluster) getNode(retrieveFn func(Storage) (string, error)) (node, error) {
	var n node
	storage := c.storage()
	address, err := retrieveFn(storage)
	if err != nil {
		return n, err
	}
	return c.getNodeByAddr(address)
}

func clientWithTimeout(dialTimeout time.Duration, fullTimeout time.Duration) *http.Client {
	transport := http.Transport{
		Dial: (&net.Dialer{
			Timeout:   dialTimeout,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: dialTimeout,
	}
	return &http.Client{
		Transport: &transport,
		Timeout:   fullTimeout,
	}
}

func (c *Cluster) getNodeByAddr(address string) (node, error) {
	var n node
	client, err := docker.NewClient(address)
	if err != nil {
		return n, err
	}
	return node{addr: address, Client: client}, nil
}
