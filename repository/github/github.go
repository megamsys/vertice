// Copyright 2015 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gandalf provides an implementation of the RepositoryManager, that
// uses Gandalf (https://github.com/tsuru/gandalf). This package doesn't expose
// any public types, in order to use it, users need to import the package and
// then configure tsuru to use the "gandalf" repo-manager.
//
package github

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	git "github.com/google/go-github/github"
	"strings"
)

func init() {
	repository.Register("github", githubManager{})
}

const endpointConfig = "git:api-server"

type githubManager struct{}

func (githubManager) client() (*gitlab.Client, error) {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: pair_token.Value},
	}
	client := git.NewClient(t.Client())
	return &client, nil
}

func (m githubManager) CreateHook(owner string, trigger string) error {
	client, err := m.client()
	if err != nil {
		return err
	}

	/*trigger_url := api_host + "/assembly/build/" + asm.Id + "/" + com.Id

	byt12 := []byte(`{"url": "","content_type": "json"}`)
	var postData map[string]interface{}
	if err := json.Unmarshal(byt12, &postData); err != nil {
		return err
	}
	postData["url"] = trigger_url

	byt1 := []byte(`{"name": "web", "active": true, "events": [ "push" ]}`)
	postHook := git.Hook{Config: postData}
	if err := json.Unmarshal(byt1, &postHook); err != nil {
		return err
	}

	pair_source, err := global.ParseKeyValuePair(com.Inputs, "source")
	if err != nil {
		log.Errorf("Failed to get the source value : %s", err.Error())
	}

	pair_owner, err := global.ParseKeyValuePair(ci.OperationRequirements, "ci-owner")
	if err != nil {
		log.Errorf("Failed to get the ci-owner value : %s", err.Error())
	}

	source := strings.Split(pair_source.Value, "/")
	_, _, err := client.Repositories.CreateHook(pair_owner.Value, strings.TrimRight(source[len(source)-1], ".git"), &postHook)

	if err != nil {
		return err
	}
	log.Debugf("[gitlab] created webhook %s successfully.", source)
	*/
	return nil

}

func (m githubManager) RemoveHook(username string) error {
	client, err := m.client()
	if err != nil {
		return err
	}
	return nil
}
