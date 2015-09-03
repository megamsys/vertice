// Copyright 2015 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gandalf provides an implementation of the RepositoryManager, that
// uses Gandalf (https://github.com/tsuru/gandalf). This package doesn't expose
// any public types, in order to use it, users need to import the package and
// then configure tsuru to use the "gandalf" repo-manager.
//
//     import _ "github.com/tsuru/tsuru/repository/gandalf"
package gitlab

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func init() {
	repository.Register("gitlab", gitlabManager{})
}

const endpointConfig = "git:api-server"

type gitlabManager struct{}

func (gitlabManager) client() (*gitlab.Client, error) {
	client := gogitlab.NewGitlab(url, version, token)
	return &client, nil
}

func (m gandalfManager) CreateHook(owner string, trigger string) error {
	client, err := m.client()
	if err != nil {
		return err
	}

	err := client.AddProjectHook(owner, trigger, false, false, false)
	if err != nil {
		return err
	}
	return nil

}

func (m gandalfManager) RemoveHook(username string) error {
	client, err := m.client()
	if err != nil {
		return err
	}
	return nil
}
