/*
** Copyright [2013-2016] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */

package one

import (
	"bytes"
	"fmt"
  "io"
	"strings"
	"text/tabwriter"

//	log "github.com/Sirupsen/logrus"
//	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/cmd"
//	"github.com/megamsys/libgo/events"
//	"github.com/megamsys/libgo/events/alerts"
//	"github.com/megamsys/libgo/utils"
//	constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/opennebula-go/api"
//	lb "github.com/megamsys/vertice/logbox"
	"github.com/megamsys/vertice/marketplaces/provision"
	"github.com/megamsys/vertice/provision/one"
	"github.com/megamsys/vertice/marketplaces/provision/one/cluster"
//	"github.com/megamsys/vertice/router"
	_ "github.com/megamsys/vertice/router/route53"
)

var mainOneProvisioner *oneProvisioner

func init() {
	mainOneProvisioner = &oneProvisioner{}
	provision.Register("one", mainOneProvisioner)
}

type oneProvisioner struct {
	defaultImage string
	vcpuThrottle string
	cluster      *cluster.Cluster
	storage       cluster.Storage
}



func (p *oneProvisioner) Cluster() *cluster.Cluster {
	if p.cluster == nil {
		panic("✗ one cluster")
	}
	return p.cluster
}

func (p *oneProvisioner) String() string {
	if p.cluster == nil {
		return "✗ one cluster"
	}
	return "ready"
}

func (p *oneProvisioner) Initialize(m interface{}) error {
	return p.initOneCluster(m)
}

func (p *oneProvisioner) initOneCluster(i interface{}) error {
	var err error
	if p.storage == nil {
		p.storage, err = buildClusterStorage()
		if err != nil {
			return err
		}
	}

	if w, ok := i.(one.One); ok {
		var nodes []cluster.Node
		p.defaultImage = w.Image
		p.vcpuThrottle = w.VCPUPercentage
		for i := 0; i < len(w.Regions); i++ {
			m := w.Regions[i].ToMap()
			c := w.Regions[i].ToClusterMap()
			n := cluster.Node{
				Address:  m[api.ENDPOINT],
				Region:   m[api.ONEZONE],
				Metadata: m,
				Clusters: c,
			}
			nodes = append(nodes, n)
		}

		//register nodes using the map.
		p.cluster, err = cluster.New(p.storage, nodes...)
		if err != nil {
			return err
		}
	}
	return nil
}

func buildClusterStorage() (cluster.Storage, error) {
	return &cluster.MapStorage{}, nil
}

// func getRouterForBox(box *provision.Box) (router.Router, error) {
// 	routerName, err := box.GetRouter()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return router.Get(routerName)
// }

func (p *oneProvisioner) StartupMessage() (string, error) {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("  > one ", "white", "", "bold") + "\t" +
		cmd.Colorfy(p.String(), "cyan", "", "")))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String()), nil
}

func (p *oneProvisioner) Create(m interface{},w io.Writer) error {
	return nil
}
