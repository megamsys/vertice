package cluster

import "time"

type Healer interface {
	HandleError(node *Node) time.Duration
}

type DefaultHealer struct{}

func (DefaultHealer) HandleError(node *Node) time.Duration {
	return 1 * time.Minute
}
