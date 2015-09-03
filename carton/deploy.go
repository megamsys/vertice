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
	"errors"
	"fmt"
	"io"
	"time"

	log "github.com/golang/glog"
	"github.com/megamsys/libgo/db"
)

type DeployOpts struct {
	B      *Box
	Config deployd.Config
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
	imageId, err := deployToProvisioner(&opts, writer)
	elapsed := time.Since(start)
	saveErr := saveDeployData(&opts, imageId, outBuffer.String(), elapsed, err)
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

func saveDeployData(opts *DeployOpts, imageId, log string, duration time.Duration, deployError error) error {
	comp := Components{
		App:       opts.Box.Name,
		Timestamp: time.Now(),
		Duration:  duration,
		Commit:    opts.Commit,
		Image:     imageId,
		Log:       log,
		User:      opts.User,
	}

	if deployError != nil {
		deploy.Error = deployError.Error()
	}

	//Riak: code to save the status of a deploy (created.)
}

func markDeploysAsRemoved(boxName string) error {
	//Riak: code to nuke a box out (component)
	return nil
}
