package repository

import (
	"fmt"
	"strings"

	"github.com/megamsys/megamd/meta"
)

const (
	//source flags
	GITHUB      = "github"
	GITLAB      = "gitlab"
	DOCKERHUB   = "dockerhub"
	MYDOCKERHUB = "mydockerhub"

	defaultManager = GITHUB

	CIHOOK   = "CI"
	TYPE     = "type"
	TOKEN    = "token"
	USERNAME = "username"

	// IMAGE indicates that the repo is an image
	IMAGE = "image"
	// Git indicates that the repo is a GIT
	GIT = "git"
	// oneclick indicates that an oneclick image exists
	ONECLICK = "oneclick"
)

var managers map[string]RepositoryManager

/* Repository represents a repository managed by the manager. */
type Repo struct {
	Type     string
	Source   string
	OneClick bool
	URL      string
	Hook     *Hook
}

type Hook struct {
	Enabled  bool
	Token    string
	UserName string
	CartonId string
	BoxId    string
}

func (r Repo) GetType() string {
	return r.Type
}

func (r Repo) GetSource() string {
	return r.Source
}

func (r Repo) Gitr() string {
	return r.URL
}

func (r Repo) IsEnabled() bool {
	return r.Hook != nil && r.Hook.Enabled
}

func (r Repo) GetToken() string {
	return r.Hook.Token
}

func (r Repo) GetUserName() string {
	return r.Hook.UserName
}

//Check on CartonId, BoxId if it exists (r.Hook.BoxId)
func (r Repo) Trigger() string {
	return meta.MC.Api + "/assembly/upgrade/" + r.Hook.CartonId
}

func (r Repo) GetShortName() (string, error) {
	i := strings.LastIndex(r.Gitr(), "/")
	if i < 0 {
		return "", fmt.Errorf("unable to parse output of git")
	}
	return strings.TrimRight(r.Gitr()[i+1:], ".git"), nil
}

//This shall be under type Tosca {} and a global method
func ForImageName(fullTosca string, version string) string {
	img := fullTosca[strings.LastIndex(fullTosca, ".")+1:]
	if len(strings.TrimSpace(version)) > 1 {
		return img + "_" + version
	}
	return img
}

type Repository interface {
	IsEnabled() bool
	GetSource() string
	GetType() string
	GetToken() string
	GetUserName() string
	Gitr() string
	Trigger() string
	GetShortName() (string, error)
}

// RepositoryManager represents a manager of application repositories.
type RepositoryManager interface {
	CreateHook(r Repository) (string, error)
	RemoveHook(r Repository) error
}

// Manager returns the current configured manager, as defined in the
// configuration file.
func Manager(managerName string) RepositoryManager {
	if _, ok := managers[managerName]; !ok {
		managerName = "nop"
	}
	return managers[managerName]
}

// Register registers a new repository manager, that can be later configured
// and used.
func Register(name string, manager RepositoryManager) {
	if managers == nil {
		managers = make(map[string]RepositoryManager)
	}
	managers[name] = manager
}
