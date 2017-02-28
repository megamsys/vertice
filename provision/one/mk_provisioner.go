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
	"fmt"
	"io"

	//	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/action"
	constants "github.com/megamsys/libgo/utils"
	lb "github.com/megamsys/vertice/logbox"
	"github.com/megamsys/vertice/provision"
)

func (p *oneProvisioner) ISODeploy(m *provision.Box, w io.Writer) error {
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
		box:           m,
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

func (p *oneProvisioner) CustomiseRawImage(m *provision.Box, w io.Writer) error {
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
		box:           m,
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
