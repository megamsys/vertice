package cluster

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/opennebula-go/metrics"
)

// Showback returns the metrics of the one cluster
func (c *Cluster) Showback(start int64, end int64) ([]interface{}, error) {
	log.Debugf("showback (%d, %d)", start, end)
	var (
		addr string
	)
	nodlist, err := c.Nodes()

	if err != nil || len(nodlist) <= 0 {
		return nil, fmt.Errorf("%s", cmd.Colorfy("Unavailable nodes (hint: start or beat it).\n", "red", "", ""))
	} else {
		addr = nodlist[0].Address
	}

	node, err := c.getNodeByAddr(addr)
	if err != nil {
		return nil, err
	}
	opts := metrics.Accounting{Api: node.Client, StartTime: start, EndTime: end}

	sb, err := opts.Get()
	if err != nil {
		return nil, wrapError(node, err)
	}
	log.Debugf("showback (%d, %d) OK", start, end)
	return sb, nil
}
