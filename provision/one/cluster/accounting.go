package cluster

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/opennebula-go/metrics"
)

// Showback returns the metrics of the one cluster
func (c *Cluster) Showback(start int64, end int64,point string) ([]interface{}, error) {
	log.Debugf("showback (%d, %d)", start, end)

	node, err := c.getNodeByAddr(point)
	if err != nil {
		return nil, fmt.Errorf("%s", cmd.Colorfy("Unavailable nodes (hint: start or beat it).\n", "red", "", ""))
	}
	opts := metrics.Accounting{Api: node.Client, StartTime: start, EndTime: end}

	sb, err := opts.Get()

	if err != nil {
		return nil, wrapError(node, err)
	}
	var a []interface{}
	a = append(a,sb)
	log.Debugf("showback (%d, %d) OK", start, end)
	return a, nil
}
