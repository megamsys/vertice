/*
** copyright [2013-2015] [Megam Systems]
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

package carton

import (
	"bytes"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megamd/provision"
)

type DeployOpts struct {
	B     *provision.Box
	Image string //why this image ?
}

// Deploy runs a deployment of an application. It will first try to run an
// image based deploy, and then fallback to the Git based deployment.
func Deploy(opts *DeployOpts) error {
	var outBuffer bytes.Buffer
	start := time.Now()
	logWriter := LogWriter{Box: opts.B}
	logWriter.Async()
	defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, &logWriter)
	imageId, err := deployToProvisioner(opts, writer)
	elapsed := time.Since(start)
	saveErr := saveDeployData(opts, imageId, outBuffer.String(), elapsed, err)
	if saveErr != nil {
		log.Errorf("WARNING: couldn't save deploy data, deploy opts: %#v", opts)
	}
	if err != nil {
		return err
	}
	return nil
}

func deployToProvisioner(opts *DeployOpts, writer io.Writer) (string, error) {
	if opts.Image != "" {
		if deployer, ok := Provisioner.(provision.ImageDeployer); ok {
			return deployer.ImageDeploy(opts.B, opts.Image, writer)
		}
	}
	return Provisioner.(provision.GitDeployer).GitDeploy(opts.B, writer)
}

func saveDeployData(opts *DeployOpts, imageId, dlog string, duration time.Duration, deployError error) error {
	log.Debugf("%s in (%s)\n%s",
		cmd.Colorfy(opts.B.GetFullName(), "cyan", "", "bold"),
		cmd.Colorfy(duration.String(), "green", "", "bold"),
		cmd.Colorfy(dlog, "yellow", "", ""))
	//if there are deployments to track as follows in outputs: {} then do it here.
	//Riak: code to save the status of a deploy (created.)
	// deploy :
	//     name:
	//     status:
	return nil
}

func markDeploysAsRemoved(boxName string) error {
	//Riak: code to nuke deploys out of a box
	return nil
}
