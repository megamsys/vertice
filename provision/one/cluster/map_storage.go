package cluster

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrNoSuchNode            = errors.New("No such node in storage")
	ErrDuplicatedNodeAddress = errors.New("Node address shouldn't repeat")
)

type MapStorage struct {
	nodes   []Node
	nodeMap map[string]*Node
	nMut    sync.Mutex
}

func (s *MapStorage) updateNodeMap() {
	s.nodeMap = make(map[string]*Node)
	for i := range s.nodes {
		s.nodeMap[s.nodes[i].Address] = &s.nodes[i]
	}
}

func (s *MapStorage) StoreNode(node Node) error {
	s.nMut.Lock()
	defer s.nMut.Unlock()
	for _, n := range s.nodes {
		if n.Address == node.Address {
			return ErrDuplicatedNodeAddress
		}
	}
	if node.Metadata == nil {
		node.Metadata = make(map[string]string)
	}
	s.nodes = append(s.nodes, node)
	s.updateNodeMap()
	return nil
}

func deepCopyNode(n Node) Node {
	newMap := map[string]string{}
	for k, v := range n.Metadata {
		newMap[k] = v
	}
	n.Metadata = newMap
	return n
}

func (s *MapStorage) RetrieveNodes() ([]Node, error) {
	s.nMut.Lock()
	defer s.nMut.Unlock()
	dst := make([]Node, len(s.nodes))
	for i := range s.nodes {
		dst[i] = deepCopyNode(s.nodes[i])
	}
	return dst, nil
}

func (s *MapStorage) RetrieveNode(address string) (Node, error) {
	s.nMut.Lock()
	defer s.nMut.Unlock()
	if s.nodeMap == nil {
		s.nodeMap = make(map[string]*Node)
	}
	node, ok := s.nodeMap[address]
	if !ok {
		return Node{}, ErrNoSuchNode
	}
	return deepCopyNode(*node), nil
}

func (s *MapStorage) UpdateNode(node Node) error {
	s.nMut.Lock()
	defer s.nMut.Unlock()
	if s.nodeMap == nil {
		s.nodeMap = make(map[string]*Node)
	}
	_, ok := s.nodeMap[node.Address]
	if !ok {
		return ErrNoSuchNode
	}
	*s.nodeMap[node.Address] = node
	return nil
}

func (s *MapStorage) RemoveNode(addr string) error {
	s.nMut.Lock()
	defer s.nMut.Unlock()
	index := -1
	for i, node := range s.nodes {
		if node.Address == addr {
			index = i
		}
	}
	if index < 0 {
		return ErrNoSuchNode
	}
	copy(s.nodes[index:], s.nodes[index+1:])
	s.nodes = s.nodes[:len(s.nodes)-1]
	s.updateNodeMap()
	return nil
}

func (s *MapStorage) RemoveNodes(addresses []string) error {
	s.nMut.Lock()
	defer s.nMut.Unlock()
	addrMap := map[string]struct{}{}
	for _, addr := range addresses {
		addrMap[addr] = struct{}{}
	}
	dup := make([]Node, 0, len(s.nodes))
	for _, node := range s.nodes {
		if _, ok := addrMap[node.Address]; !ok {
			dup = append(dup, node)
		}
	}
	if len(dup) == len(s.nodes) {
		return ErrNoSuchNode
	}
	s.nodes = dup
	s.updateNodeMap()
	return nil
}

func (s *MapStorage) LockNodeForHealing(address string, isFailure bool, timeout time.Duration) (bool, error) {
	s.nMut.Lock()
	defer s.nMut.Unlock()
	n, present := s.nodeMap[address]
	if !present {
		return false, ErrNoSuchNode
	}
	now := time.Now().UTC()
	if n.Healing.LockedUntil.After(now) {
		return false, nil
	}
	n.Healing.LockedUntil = now.Add(timeout)
	n.Healing.IsFailure = isFailure
	return true, nil
}

func (s *MapStorage) ExtendNodeLock(address string, timeout time.Duration) error {
	s.nMut.Lock()
	defer s.nMut.Unlock()
	n, present := s.nodeMap[address]
	if !present {
		return ErrNoSuchNode
	}
	now := time.Now().UTC()
	n.Healing.LockedUntil = now.Add(timeout)
	return nil
}

func (s *MapStorage) UnlockNode(address string) error {
	s.nMut.Lock()
	defer s.nMut.Unlock()
	n, present := s.nodeMap[address]
	if !present {
		return ErrNoSuchNode
	}
	n.Healing = HealingData{}
	return nil
}
