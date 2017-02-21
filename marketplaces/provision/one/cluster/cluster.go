package cluster

import (
	"errors"
	"fmt"
	"github.com/megamsys/libgo/utils"
	"github.com/megamsys/opennebula-go/api"
	"sync/atomic"
	"time"
)

var (
	errStorageMandatory = errors.New("Storage parameter is mandatory")
	errHealerInProgress = errors.New("Healer already running")
)

type node struct {
	addr     string
	template string
	image    string
	Client   *api.Rpc
}

type NodeStorage interface {
	StoreNode(node Node) error
	RetrieveNodes() ([]Node, error)
	RetrieveNode(address string) (Node, error)
	UpdateNode(node Node) error
	RemoveNode(address string) error
	RemoveNodes(addresses []string) error
	LockNodeForHealing(address string, isFailure bool, timeout time.Duration) (bool, error)
	ExtendNodeLock(address string, timeout time.Duration) error
	UnlockNode(address string) error
}

type Storage interface {
	NodeStorage
}

type ClusterHook interface {
	BeforeCreateMachine(node Node) error
}

// Cluster is the basic type of the package. It manages internal nodes, and
// provide methods for interaction with those nodes
type Cluster struct {
	Healer Healer
	Hook   ClusterHook
	stor   Storage
}

type OneNodeError struct {
	node node
	cmd  string
	err  error
}

func (n OneNodeError) Error() string {
	if n.cmd == "" {
		return fmt.Sprintf("error in one node %q: %s", n.node.addr, n.err.Error())
	}
	return fmt.Sprintf("error in one node %q running command %q: %s", n.node.addr, n.cmd, n.err.Error())
}

func (n OneNodeError) BaseError() error {
	return n.err
}

func wrapError(n node, err error) error {
	if err != nil {
		return OneNodeError{node: n, err: err}
	}
	return nil
}

func wrapErrorWithCmd(n node, err error, cmd string) error {
	if err != nil {
		return OneNodeError{node: n, err: err, cmd: cmd}
	}
	return nil
}

// New creates a new Cluster, initially composed by the given nodes.
// The storage parameter is the storage the cluster instance will use.
func New(storage Storage, nodes ...Node) (*Cluster, error) {
	var (
		c   Cluster
		err error
	)
	if storage == nil {
		return nil, errStorageMandatory
	}
	c.stor = storage
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
	if node.Region == "" {
		return errors.New("Invalid region")
	}
	return c.storage().StoreNode(node)
}

func (c *Cluster) UpdateNode(node Node) (Node, error) {
	_, err := c.storage().RetrieveNode(node.Region)
	if err != nil {
		return Node{}, err
	}
	unlock, err := c.lockWithTimeout(node.Region, false)
	if err != nil {
		return Node{}, err
	}
	defer unlock()
	dbNode, err := c.storage().RetrieveNode(node.Region)
	if err != nil {
		return Node{}, err
	}
	if node.CreationStatus != "" && node.CreationStatus != dbNode.CreationStatus {
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
func (c *Cluster) Unregister(region string) error {
	return c.storage().RemoveNode(region)
}

func (c *Cluster) UnregisterNodes(regions ...string) error {
	return c.storage().RemoveNodes(regions)
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

func (c *Cluster) handleNodeError(region string, lastErr error, incrementFailures bool) error {
	unlock, err := c.lockWithTimeout(region, true)
	if err != nil {
		return err
	}
	go func() {
		defer unlock()
		node, err := c.storage().RetrieveNode(region)
		if err != nil {
			return
		}
		node.updateError(lastErr, incrementFailures)
		duration := c.Healer.HandleError(&node)
		if duration > 0 {
			node.updateDisabled(time.Now().Add(duration))
		}
		c.storage().UpdateNode(node)
		if fn := nodeUpdatedOnError.Val(); fn != nil {
			fn()
		}
	}()
	return nil
}

// Modified by tests
var nodeUpdatedOnError nodeUpdatedHook

type nodeUpdatedHook struct {
	atomic.Value
}

func (v *nodeUpdatedHook) Val() func() {
	if fn := v.Load(); fn != nil {
		return fn.(func())
	}
	return nil
}

func (c *Cluster) handleNodeSuccess(region string) error {
	unlock, err := c.lockWithTimeout(region, false)
	if err != nil {
		return err
	}
	defer unlock()
	node, err := c.storage().RetrieveNode(region)
	if err != nil {
		return err
	}
	node.updateSuccess()
	return c.storage().UpdateNode(node)
}

func (c *Cluster) lockWithTimeout(region string, isFailure bool) (func(), error) {
	lockTimeout := 3 * time.Minute
	locked, err := c.storage().LockNodeForHealing(region, isFailure, lockTimeout)
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
			c.storage().ExtendNodeLock(region, lockTimeout)
		}
	}()
	return func() {
		doneKeepAlive <- true
		c.storage().UnlockNode(region)
	}, nil
}

func (c *Cluster) storage() Storage {
	return c.stor
}

func (c *Cluster) getNode(retrieveFn func(Storage) (Node, error)) (node, error) {
	var n node
	storage := c.storage()
	node, err := retrieveFn(storage)
	if err != nil {
		return n, err
	}
	return c.getNodeByObject(node)
}

func (c *Cluster) getNodeByObject(nodeo Node) (node, error) {
	var n node
	client, err := api.NewClient(map[string]string{api.ENDPOINT: nodeo.Address, api.USERID: nodeo.Metadata[api.USERID], api.PASSWORD: nodeo.Metadata[api.PASSWORD]})

	if err != nil {
		return n, err
	}

	template := nodeo.Metadata[api.TEMPLATE]
	return node{addr: nodeo.Address, template: template, Client: client}, nil
}

//return vnets and cluster id which is choosen
func (c *Cluster) getVnets(nodeo Node, m map[string]string, st string) (map[string]string, string) {
	res := make(map[string]string)
	for k, v := range nodeo.Clusters {
		if v[utils.STORAGE_TYPE] == st && v[utils.VONE_CLOUD] != utils.TRUE {
			for i, j := range nodeo.Clusters[k] {
				if m[i] == utils.TRUE {
					res[i] = j
				}
			}
			return res, k
		}
	}
	return res, ""
}
