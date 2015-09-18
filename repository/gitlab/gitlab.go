package gitlab

import (
	"github.com/megamsys/megamd/repository"
	"github.com/plouc/go-gitlab-client"
)

func init() {
	repository.Register("gitlab", gitlabManager{})
}

const endpointConfig = "git:api-server"

type gitlabManager struct{}

func (gitlabManager) client() (*gogitlab.Gitlab, error) {
	url, version, token := "", "", ""
	return gogitlab.NewGitlab(url, version, token), nil
}

func (m gitlabManager) CreateHook(owner string, trigger string) error {
	client, err := m.client()
	if err != nil {
		return err
	}

	err = client.AddProjectHook(owner, trigger, false, false, false)
	if err != nil {
		return err
	}
	return nil

}

func (m gitlabManager) RemoveHook(owner string) error {
	/*client, err := m.client()
	if err != nil {
		return err
	}*/
	return nil
}
