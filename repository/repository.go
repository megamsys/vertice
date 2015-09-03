// Copyright 2015 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package repository contains types and functions for git repository
// interaction.
package repository

import (
	"errors"

)

const defaultManager = "github"

var managers map[string]RepositoryManager

// Repository represents a repository in the manager.
type Repository struct {
	Name         string
	ReadOnlyURL  string
	ReadWriteURL string
}

// RepositoryManager represents a manager of application repositories.
type RepositoryManager interface {
	CreateHook(username string) error
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
