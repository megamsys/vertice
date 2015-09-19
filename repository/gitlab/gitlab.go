package gitlab

import (
	"github.com/megamsys/megamd/repository"
	"github.com/plouc/go-gitlab-client"
)

func init() {
	repository.Register("gitlab", gitlabManager{})
}

type gitlabManager struct{}

//http://base_url/api_path/projects?private_token=token")
func (gitlabManager) client() (*gogitlab.Gitlab, error) {
	url, version, token := "http://base_url", "api_path", "token"
	return gogitlab.NewGitlab(url, version, token), nil
}

func (m gitlabManager) CreateHook(r repository.Repository) (string, error) {
	client, err := m.client()
	if err != nil {
		return "", err
	}
	err = client.AddProjectHook(r.GetUserName(), r.Trigger(), false, false, false)
	if err != nil {
		return "", err
	}
	return "", nil

}

func (m gitlabManager) RemoveHook(r repository.Repository) error {
	/*client, err := m.client()
	if err != nil {
		return err
	}*/
	return nil
}
