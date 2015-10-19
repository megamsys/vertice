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
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megamd/provision"
	"github.com/megamsys/megamd/repository"
)

type DeployData struct {
	BoxName     string
	HookId      string
	PrivateIp   string
	PublicIp    string
	Timestamp   time.Time
	Duration    time.Duration
	Commit      string
	Image       string
	Origin      string
	CanRollback bool
}

type DeployOpts struct {
	B *provision.Box
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
	if opts.B.Repo == nil || opts.B.Repo.Type == repository.IMAGE {
		if deployer, ok := Provisioner.(provision.ImageDeployer); ok {
			return deployer.ImageDeploy(opts.B, image(opts.B), writer)
		}
	}
	return Provisioner.(provision.GitDeployer).GitDeploy(opts.B, writer)
}

// for a vm provisioner return the last name (tosca.torpedo.ubuntu) ubuntu as the image name.
// for docker return the Inputs[image]
func image(b *provision.Box) string {
	if b.Repo == nil {
		img := b.Tosca[strings.LastIndex(b.Tosca, ".")+1:]
		if len(strings.TrimSpace(b.ImageVersion)) > 1 {
			return img + "_" + b.ImageVersion
		}
		return img
	}
	return ""
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

	/*deploy := DeployData {
		App:       opts.App.Name,
		Timestamp: time.Now(),
		Duration:  duration,
		Commit:    opts.Commit,
		Image:     imageId,
		Log:       log,
	}
	if opts.Commit != "" {
		deploy.Origin = "git"
	} else if opts.Image != "" {
		deploy.Origin = "rollback"
	} else {
		deploy.Origin = "app-deploy"
	}
	if deployError != nil {
		deploy.Error = deployError.Error()
	}
	return db.Store(compid or assmid, &struct)
	*/
	return nil
}
