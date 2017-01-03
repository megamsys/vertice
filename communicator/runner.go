package communicator


import (
	"fmt"
)
var managers map[string]Runner

// Runner  interface.
type Runner interface {

	// Stringified version of runner
	String() string

	Rerun() (r Runner, err error)

	CleanUp() ( err error)
	// Storage operations.
	Run(packages []string) (r Runner, err error)
}

type BaseHost struct {
	Host string  `json:"ipaddress"`
	Port string  `json:"port"`
	Username string `json:"username"`
	SSHKey
	SSHPassword
	IsRan bool    //to verify can cleanup or not
}

type SSHKey struct {
	PrivateKey string `json:"privatekey"`
	PublicKey  string `json:"publickey"`
}

type SSHPassword struct {
  Password string `json:"password"`
}

type Run interface {
}

func Get(name string) (Run, error) {
	p, ok := managers[name]
	if !ok {
		return nil, fmt.Errorf("unknown runner : %q", name)
	}
	return p, nil
}

// Manager returns the current configured manager, as defined in the
// configuration file.
func Manager(managerName string) Runner {
	if _, ok := managers[managerName]; !ok {
		managerName = "nop"
	}
	return managers[managerName]
}

// Register registers a new Runner manager, that can be later configured
// and used.
func Register(name string, manager Runner) {
	if managers == nil {
		managers = make(map[string]Runner)
	}
	managers[name] = manager
}
