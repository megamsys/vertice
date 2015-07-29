/*
** Copyright [2013-2015] [Megam Systems]
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
package provisioner

import (
	"fmt"
	"github.com/megamsys/megamd/global"
)

/*
 * Provisioner is the basic interface of this package.
 */
type Provisioner interface {

  // Provision is called when megam engine is creating the app.
	Create(*global.AssemblyWithComponents, string, bool, string) (string, error)

	// Delete is called when megam is destroying the app.
	Delete(*global.AssemblyWithComponents, string) (string, error)

	// Start starts the units of the application, with an optional string
	// parameter represeting the name of the process to start. When the
	// process is empty, Start will start all units of the application.
	Start(App, string) error

	// Stop stops the units of the application, with an optional string
	// parameter represeting the name of the process to start. When the
	// process is empty, Stop will stop all units of the application.
	Stop(App, string) error

	// Restart restarts the units of the application, with an optional
	// string parameter represeting the name of the process to start. When
	// the process is empty, Restart will restart all units of the
	// application.
	Restart(App, string, io.Writer) error

	// SetUnitStatus changes the status of a unit.
	SetUnitStatus(Unit, Status) error

	// ExecuteCommand runs a command in all units of the app.
	ExecuteCommand(stdout, stderr io.Writer, app App, cmd string, args ...string) error

	// Addr returns the address for an app.
	//
	// megam will use this method to get the IP (although it might not be
	// an actual IP, collector calls it "IP") of the app from the
	// provisioner.
	Addr(App) (string, error)

	// Register a unit after the container has been created or restarted.
	RegisterUnit(Unit, map[string]interface{}) error

	// Open a remote shel in one of the units in the application.
	Shell(ShellOptions) error

	// Returns the metric backend environs for the app.
	MetricEnvs(App) map[string]string


}

// CNameManager represents a provisioner that supports cname on applications.
type CNameManager interface {
	SetCName(app App, cname string) error
	UnsetCName(app App, cname string) error
}

// ShellOptions is the set of options that can be used when calling the method
// Shell in the provisioner.
type ShellOptions struct {
	App    App
	Conn   io.ReadWriteCloser
	Width  int
	Height int
	Unit   string
	Term   string
}


// GitDeployer is a provisioner that can deploy the application from a Git
// repository.
type GitDeployer interface {
	GitDeploy(app App, version string, w io.Writer) (string, error)
}


// ImageDeployer is a provisioner that can deploy the application from a
// previously generated image.
type ImageDeployer interface {
	ImageDeploy(app App, image string, w io.Writer) (string, error)
}


type MessageProvisioner interface {
	StartupMessage() (string, error)
}

// InitializableProvisioner is a provisioner that provides an initialization
// method that should be called when the engine is started
type InitializableProvisioner interface {
	Initialize() error
}


var provisioners = make(map[string]Provisioner)
/*
 * Register registers a new provisioner in the Provisioner registry.
 */
func Register(name string, p Provisioner) {
	provisioners[name] = p
}

func GetProvisioner(name string) (Provisioner, error) {
	provider, ok := provisioners[name]
	if !ok {
		return nil, fmt.Errorf("Provisioner not registered")
	}
	return provider, nil
}
