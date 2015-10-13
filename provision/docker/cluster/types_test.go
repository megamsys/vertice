package cluster

import (
	"errors"
	"time"

	"github.com/fsouza/go-dockerclient"
)

type containerList []docker.APIContainers

func (l containerList) Len() int {
	return len(l)
}

func (l containerList) Less(i, j int) bool {
	return l[i].ID < l[j].ID
}

func (l containerList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type failingStorage struct{}

func (failingStorage) StoreContainer(container, host string) error {
	return errors.New("storage error")
}
func (failingStorage) RetrieveContainer(container string) (string, error) {
	return "", errors.New("storage error")
}
func (failingStorage) RemoveContainer(container string) error {
	return errors.New("storage error")
}
func (failingStorage) RetrieveContainers() ([]Container, error) {
	return nil, errors.New("storage error")
}
func (failingStorage) StoreImage(repository, id, host string) error {
	return errors.New("storage error")
}
func (failingStorage) RetrieveImage(repository string) (Image, error) {
	return Image{}, errors.New("storage error")
}
func (failingStorage) RemoveImage(repository, id, host string) error {
	return errors.New("storage error")
}
func (failingStorage) RetrieveImages() ([]Image, error) {
	return nil, errors.New("storage error")
}
func (failingStorage) StoreNode(node Node) error {
	return errors.New("storage error")
}
func (failingStorage) RetrieveNodesByMetadata(metadata map[string]string) ([]Node, error) {
	return nil, errors.New("storage error")
}
func (failingStorage) RetrieveNodes() ([]Node, error) {
	return nil, errors.New("storage error")
}
func (failingStorage) RetrieveNode(addr string) (Node, error) {
	return Node{}, errors.New("storage error")
}
func (failingStorage) UpdateNode(node Node) error {
	return errors.New("storage error")
}
func (failingStorage) RemoveNode(address string) error {
	return errors.New("storage error")
}
func (failingStorage) LockNodeForHealing(address string, isFailure bool, timeout time.Duration) (bool, error) {
	return false, errors.New("storage error")
}
func (failingStorage) ExtendNodeLock(address string, timeout time.Duration) error {
	return errors.New("storage error")
}
func (failingStorage) UnlockNode(address string) error {
	return errors.New("storage error")
}
