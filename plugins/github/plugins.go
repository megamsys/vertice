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
package github

import (
	"golang.org/x/oauth2"
	"encoding/json"
	git "github.com/google/go-github/github"
	"github.com/megamsys/megamd/global"
	"github.com/megamsys/megamd/log"
	"github.com/megamsys/megamd/plugins"
	"strings"
)

/**
**Github register function
**This function register to plugins container
**/
func Init() {
	plugins.RegisterPlugins("github", &GithubPlugin{})
}

type GithubPlugin struct{}

const (
	GITHUB = "github"
	ENABLE = "true"
	CI     = "CI"
)

/**
**watcher function executes the github repository webhook creation
**first get the ci value and parse it and to create the hook for that users repository
**get trigger url from config file
**/
func (c *GithubPlugin) Watcher(asm *app.DeepAssembly, ci *global.Operations, com *global.Component) error {
	switch ci.OperationType {
	case CI:
		cierr := cioperation(asm, ci, com)
		if cierr != nil {
			return cierr
		}
		break
	}
	return nil
}

func cioperation(asm *app.DeepAssembly, ci *global.Operations, com *global.Component) error {
	pair_scm, err := global.ParseKeyValuePair(ci.OperationRequirements, "ci-scm")
	if err != nil {
		log.Errorf("Failed to get the ci-scm value : %s", err.Error())
	}

	pair_enable, err := global.ParseKeyValuePair(ci.OperationRequirements, "ci-enable")
	if err != nil {
		log.Errorf("Failed to get the ci-enable value : %s", err.Error())
	}

	if pair_scm.Value == GITHUB && pair_enable.Value == ENABLE {
		pair_token, err := global.ParseKeyValuePair(ci.OperationRequirements, "ci-token")
		if err != nil {
			log.Errorf("Failed to get the ci-token value : %s", err.Error())
		}
		t := &oauth.Transport{
			Token: &oauth.Token{AccessToken: pair_token.Value},
		}

		client := git.NewClient(t.Client())

		api_host, err := config.GetString("megam:api")
		if err != nil {
			return err
		}

		trigger_url := api_host + "/assembly/build/" + asm.Id + "/" + com.Id

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

	} else {
		log.Debugf("[github] skipped...")
	}
	return nil
}

/**
**notify the messages or any other operations from github
** now this function performs build the pushed application from github to remote
**/
func (c *GithubPlugin) Notify(m *global.EventMessage) error {
	return nil
}
