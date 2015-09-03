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
package provision

import (
	"errors"
	"fmt"
	"io"
	"net/url"

)

var (
	ErrInvalidStatus = errors.New("invalid status")
	ErrEmptyApp      = errors.New("no units for this app")
	ErrUnitNotFound  = errors.New("unit not found")
)

// Status represents the status of a unit in megamd
type Status string

func (s Status) String() string {
	return string(s)
}

func ParseStatus(status string) (Status, error) {
	switch status {
	case "created":
		return StatusCreated, nil
	case "building":
		return StatusBuilding, nil
	case "error":
		return StatusError, nil
	case "started":
		return StatusStarted, nil
	case "starting":
		return StatusStarting, nil
	case "stopped":
		return StatusStopped, nil
	}
	return Status(""), ErrInvalidStatus
}

const (
	// StatusCreated is the initial status of a unit in the database,
	// it should transition shortly to a more specific status
	StatusCreated = Status("created")

	// StatusBuilding is the status for units being provisioned by the
	// provisioner, like in the deployment.
	StatusBuilding = Status("building")

	// StatusError is the status for units that failed to start, because of
	// an application error.
	StatusError = Status("error")

	// StatusStarting is set when the container is started in docker.
	StatusStarting = Status("starting")

	// StatusStarted is for cases where the unit is up and running, and bound
	// to the proper status, it's set by RegisterUnit and SetUnitStatus.
	StatusStarted = Status("started")

	// StatusStopped is for cases where the unit has been stopped.
	StatusStopped = Status("stopped")
)

// Box represents a provision unit. Can be a machine, container or anything
// IP-addressable.
type Box struct {
	C          carton.Component
	Name       string
	DomainName string
	Tosca      string
	Commit     string
	Image      string
	Git        string
	Status     Status
	Provider   string
	Address    *url.URL
	Ip         string
}

// GetName returns the name of the box.
func (b *Box) GetFullName() string {
	return b.Name + b.DomainName
}

// GetTosca returns the tosca type of the box.
func (b *Box) GetTosca() string {
	return b.Tosca
}

// GetIp returns the Unit.IP.
func (b *Box) GetIp() string {
	return b.Ip
}

// Available returns true if the unit is available. It will return true
// whenever the unit itself is available, even when the application process is
// not.
func (b *Box) Available() bool {
	return b.Status == StatusStarted ||
		b.Status == StatusStarting ||
		b.Status == StatusError
}

// Named is something that has a name, providing the GetName method.
type Named interface {
	GetName() string
}

// App represents a megamd app.
//
// It contains only relevant information for provisioning.
type Carton interface {
	Named

	Bind(*Box) error
	Unbind(*Unit) error

	// Log should be used to log messages in the app.
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

// CNameManager represents a provisioner that supports cname on applications.
type CNameManager interface {
	SetCName(app App, cname string) error
	UnsetCName(app App, cname string) error
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

// Provisioner is the basic interface of this package.
//
// Any megamd provisioner must implement this interface in order to provision
// megamd apps.
type Provisioner interface {
	// Destroy is called when megamd is destroying the box.
	Destroy(*Box, io.Writer) error

	// SetBoxStatus changes the status of a unit.
	SetBoxStatus(*Box, Status) error

	// ExecuteCommand runs a command in all boxes of the carton.
	ExecuteCommand(stdout, stderr io.Writer, c Carton, cmd string, args ...string) error

	// ExecuteCommandOnce runs a command in one box of the carton.
	ExecuteCommandOnce(stdout, stderr io.Writer, c Carton, cmd string, args ...string) error

	// Restart restarts the units of the application, with an optional
	// string parameter represeting the name of the process to start. When
	// the process is empty, Restart will restart all units of the
	// application.
	Restart(*Carton, string, io.Writer) error

	// Start starts the units of the application, with an optional string
	// parameter represeting the name of the process to start. When the
	// process is empty, Start will start all units of the application.
	Start(*Carton, string) error

	// Stop stops the units of the application, with an optional string
	// parameter represeting the name of the process to start. When the
	// process is empty, Stop will stop all units of the application.
	Stop(*Carton, string) error

	// Addr returns the address for an app.
	//
	// megamd will use this method to get the IP (although it might not be
	// an actual IP, collector calls it "IP") of the app from the
	// provisioner.
	Addr(Carton) (string, error)

	// Returns the metric backend environs for the app.
	MetricEnvs(App) map[string]string
}

type MessageProvisioner interface {
	StartupMessage() (string, error)
}

// InitializableProvisioner is a provisioner that provides an initialization
// method that should be called when the app is started
type InitializableProvisioner interface {
	Initialize() error
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

type MegamdYamlRestartHooks struct {
	Before []string
	After  []string
}

type MegamdYamlHooks struct {
	Restart MegamdYamlRestartHooks
	Build   []string
}

type MegamdYamlHealthcheck struct {
	Path            string
	Method          string
	Status          int
	Match           string
	AllowedFailures int `json:"allowed_failures" bson:"allowed_failures"`
}

type MegamdYamlData struct {
	Hooks       MegamdYamlHooks
	Healthcheck MegamdYamlHealthcheck
}
