// Copyright 2015 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package repository contains types and functions for git repository
// interaction.
package repository

const (
	defaultManager = "github"
	CI             = "CI"
	CI_ENABLED     = "ci-enabled"
	CI_TOKEN       = "ci-token"
	CI_SCM         = "ci-scm"
	CI_USER        = "ci-user"
	CI_URL         = "ci-url"
	CI_APIVERSION  = "ci-apiversion"
)

var managers map[string]RepositoryManager

/* Repository represents a repository in the manager. */
type Repo struct {
	Enabled  bool
	Token    string
	Git      string
	GitURL   string
	UserName string
	Version  string
}

func (r Repo) IsEnabled() bool {
	return r.Enabled
}

func (r Repo) GetToken() string {
	return r.Token
}

func (r Repo) GetGit() string {
	return r.Git
}

func (r Repo) GetGitURL() string {
	return r.GitURL
}

func (r Repo) GetUserName() string {
	return r.UserName
}

func (r Repo) GetVersion() string {
	return r.Version
}

type Repository interface {
	IsEnabled() bool
	GetToken() string
	//GetGit() string
	//GetGitURL() string
	GetUserName() string
	GetVersion() string
}

// RepositoryManager represents a manager of application repositories.
type RepositoryManager interface {
	CreateHook(username string, trigger string) error
	RemoveHook(username string) error
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
