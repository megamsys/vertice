package github

import (
	//log "github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
	"github.com/megamsys/megamd/repository"
	"golang.org/x/oauth2"
)

func init() {
	repository.Register("github", githubManager{})
}

const endpointConfig = "git:api-server"

type githubManager struct{}

func (githubManager) client() (*github.Client, error) {
	token := ""

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	return github.NewClient(oauth2.NewClient(oauth2.NoContext, ts)), nil //there is no error trap here ?

}

func (m githubManager) CreateHook(owner string, trigger string) error {
	/*client, err := m.client()
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
	*/
	return nil

}

func (m githubManager) RemoveHook(username string) error {
	/*client, err := m.client()
	if err != nil {
		return err
	}
	*/
	return nil
}
