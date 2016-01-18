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
package provision

import (
	"errors"
	"fmt"
	"io"

	"github.com/megamsys/megamd/carton/bind"
	"github.com/megamsys/megamd/repository"
)

const (
	PROVIDER        = "provider"
	PROVIDER_ONE    = "one"
	PROVIDER_DOCKER = "docker"
)

var (
	ErrInvalidStatus  = errors.New("invalid status")
	ErrEmptyCarton    = errors.New("no boxs for this carton")
	ErrBoxNotFound    = errors.New("box not found")
	ErrNoOutputsFound = errors.New("no outputs found in the box. Did you set it ? ")
	ErrNotImplemented = errors.New("I'am on diet.")
)

// Status represents the status of a unit in megamd
type Status string

func (s Status) String() string {
	return string(s)
}

const (
	// StatusLaunching is the initial status of a box
	// it should transition shortly to a more specific status
	StatusLaunching = Status("launching")

	// StatusLaunched is the status for box after launched in cloud.
	StatusLaunched = Status("launched")

	// StatusBootstrapped is the status for box after being booted by the agent in cloud
	StatusBootstrapped = Status("boostrapped")

	// Stateup is the status of the which is moving up in the state in cloud.
	// Sent by megamd to gulpd when it received StatusBootstrapped.
	StatusStateup = Status("stateup")

	//fully up instance
	StatusRunning = Status("running")

	StatusStarting = Status("starting")
	StatusStarted  = Status("started")

	StatusStopping = Status("stopping")
	StatusStopped  = Status("stopped")

	StatusUpgraded = Status("upgraded")

	StatusDestroying = Status("destroying")
	StatusNuked      = Status("nuked")

	// StatusError is the status for units that failed to start, because of
	// a box error.
	StatusError = Status("error")
)

var Repositorys repository.Repositorys

// Named is something that has a name, providing the GetName method.
type Named interface {
	GetName() string
}

// Carton represents a deployment entity in megamd.
//
// It contains boxes to provision and only relevant information for provisioning.
type Carton interface {
	Named

	Bind(*Box) error
	Unbind(*Box) error

	// Log should be used to log messages in the box.
	Log(message, source, unit string) error

	Boxes() []*Box

	// Run executes the command in box units. Commands executed with this
	// method should have access to environment variables defined in the
	// app.
	Run(cmd string, w io.Writer, once bool) error

	Envs() map[string]bind.EnvVar

	GetMemory() int64
	GetSwap() int64
	GetCpuShare() int
}

// CNameManager represents a provisioner that supports cname on box.
type CNameManager interface {
	SetCName(b *Box, cname string) error
	UnsetCName(b *Box, cname string) error
}

// ShellOptions is the set of options that can be used when calling the method
// Shell in the provisioner.
type ShellOptions struct {
	Box    *Box
	Conn   io.ReadWriteCloser
	Width  int
	Height int
	Unit   string
	Term   string
}

// GitDeployer is a provisioner that can deploy the box from a Git
// repository.
type GitDeployer interface {
	GitDeploy(b *Box, w io.Writer) (string, error)
}

// ImageDeployer is a provisioner that can deploy the box from a
// previously generated image.
type ImageDeployer interface {
	ImageDeploy(b *Box, image string, w io.Writer) (string, error)
}

// StateChanger changes the state of a deployed box
// A deployed box is termed as a machine or a container
type StateChanger interface {
	SetState(*Box, io.Writer, Status) error
}

// Provisioner is the basic interface of this package.
//
// Any megamd provisioner must implement this interface in order to provision
// megamd cartons.
type Provisioner interface {

	// Destroy is called when megamd is destroying the box.
	Destroy(*Box, io.Writer) error

	// SetBoxStatus changes the status of a box.
	SetBoxStatus(*Box, io.Writer, Status) error

	// ExecuteCommandOnce runs a command in one box of the carton.
	ExecuteCommandOnce(stdout, stderr io.Writer, box *Box, cmd string, args ...string) error

	// Restart restarts the boxes of the carton, with an optional
	// string parameter represeting the name of the process to start.
	Restart(*Box, string, io.Writer) error
	// Start starts the boxes of the application, with an optional string
	// parameter represeting the name of the process to start.
	Start(*Box, string, io.Writer) error

	// Stop stops the boxes of the application, with an optional string
	// parameter represeting the name of the process to stop.
	Stop(*Box, string, io.Writer) error

	// Open a remote shel in one of the boxs in the carton.
	Shell(ShellOptions) error

	// Addr returns the address for an box.
	//
	// megamd will use this method to get the IP (although it might not be
	// an actual IP, collector calls it "IP") of the app from the
	// provisioner.
	Addr(*Box) (string, error)

	// Returns the metric backend environs for the npx.
	MetricEnvs(Carton) map[string]string
}

type MessageProvisioner interface {
	StartupMessage() (string, error)
}

// InitializableProvisioner is a provisioner that provides an initialization
// method that should be called when the carton is started,
//additionally provide a map of configuration info.
type InitializableProvisioner interface {
	Initialize(m map[string]string, b map[string]string) error
}

// ExtensibleProvisioner is a provisioner where administrators can manage
// platforms (automatically adding, removing and updating platforms).
type ExtensibleProvisioner interface {
	PlatformAdd(name string, args map[string]string, w io.Writer) error
	PlatformUpdate(name string, args map[string]string, w io.Writer) error
	PlatformRemove(name string) error
}

var provisioners = make(map[string]Provisioner)

// Register registers a new provisioner in the Provisioner registry.
func Register(name string, p Provisioner) {
	provisioners[name] = p
}

// Get gets the named provisioner from the registry.
func Get(name string) (Provisioner, error) {
	p, ok := provisioners[name]
	if !ok {
		return nil, fmt.Errorf("unknown provisioner: %q", name)
	}
	return p, nil
}

// Registry returns the list of registered provisioners.
func Registry() []Provisioner {
	registry := make([]Provisioner, 0, len(provisioners))
	for _, p := range provisioners {
		registry = append(registry, p)
	}
	return registry
}

// Error represents a provisioning error. It encapsulates further errors.
type Error struct {
	Reason string
	Err    error
}

// Error is the string representation of a provisioning error.
func (e *Error) Error() string {
	var err string
	if e.Err != nil {
		err = e.Err.Error() + ": " + e.Reason
	} else {
		err = e.Reason
	}
	return err
}
