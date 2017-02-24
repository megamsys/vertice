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
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/cmd"
	//	"github.com/megamsys/libgo/events"
	//	"github.com/megamsys/libgo/events/alerts"
	//	"github.com/megamsys/libgo/utils"
	constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/opennebula-go/api"
	lb "github.com/megamsys/vertice/logbox"
	"github.com/megamsys/vertice/marketplaces/provision"
	"github.com/megamsys/vertice/provision/one"
	//	"github.com/megamsys/vertice/marketplaces"
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
	vcpu         string
	disksize     string
	ram          string
	cluster      *cluster.Cluster
	storage      cluster.Storage
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

func (p *oneProvisioner) Initialize(m interface{}, mkc map[string]string) error {
	return p.initOneCluster(m, mkc)
}

func (p *oneProvisioner) initOneCluster(i interface{}, mkc map[string]string) error {
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
		p.vcpu = mkc[constants.CPU]
		p.ram = mkc[constants.RAM]
		p.disksize = mkc[constants.STORAGE]

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

func (p *oneProvisioner) ISODeploy(box interface{}, w io.Writer) error {
  m := box.(provision.Box)
	fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("--- deploy box (%s)", m.Name)))

	actions := []*action.Action{
		&machCreating,
		&createRawISOImage,
		&updateRawImageId,
		&waitUntillImageReady,
		&updateRawStatus,
	}
	pipeline := action.NewPipeline(actions...)
	args := runMachineActionsArgs{
		box:           &m,
		writer:        w,
		machineStatus: constants.StatusCreating,
		provisioner:   p,
	}

	err := pipeline.Execute(args)
	if err != nil {
		fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("--- create iso pipeline for box (%s)\n --> %s", m.Name, err)))
		return err
	}

	fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("--- create iso pipeline for box (%s)OK", m.Name)))
	return nil
}

func (p *oneProvisioner) CustomiseRawImage(box interface{}, w io.Writer) error {
  m := box.(provision.Box)
	fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("--- customize rawimage pipeline for box (%s)", m.Name)))

	actions := []*action.Action{
		&machCreating,
	  &createDatablockImage,
		&updateMarketplaceStatus,
	  &updateMarketplaceImageId,
	  &waitUntillImageReady,
		&updateMarketplaceStatus,
		// &createInstanceForCustomize,
		// &updateMarketplaceStatus,
		// &waitUntillvmReady,
		// &updateVncHostIp,
	  // &updateMarketplaceStatus,
	}
	pipeline := action.NewPipeline(actions...)
	args := runMachineActionsArgs{
		box:           &m,
		writer:        w,
		machineStatus: constants.StatusCreating,
		provisioner:   p,
	}

	err := pipeline.Execute(args)
	if err != nil {
		fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.ERROR, fmt.Sprintf("--- Error:  customize rawimage pipeline for box (%s)\n --> %s", m.Name, err)))
		return err
	}

	fmt.Fprintf(w, lb.W(lb.DEPLOY, lb.INFO, fmt.Sprintf("--- customize rawimage pipeline for box (%s)OK", m.Name)))
	return nil
}

func (o *oneProvisioner) Resource() map[string]string {
	m := make(map[string]string)
	m[constants.CPU] = o.vcpu
	m[constants.RAM] = o.ram
	m[constants.STORAGE] = o.disksize
	return m
}
