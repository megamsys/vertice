package github

import (
  "strconv"
  "net/http"
  "net/http/httputil"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megamd/repository"
	"golang.org/x/oauth2"
)

func init() {
	repository.Register("github", githubManager{})
}

type githubManager struct{}

func (m githubManager) client(token string) *github.Client {
	return github.NewClient(oauth2.NewClient(oauth2.NoContext,
		oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)))
}

//https://developer.github.com/v3/repos/hooks/#create-a-hook
func (m githubManager) CreateHook(r repository.Repository) (string, error) {
	client := m.client(r.GetToken())
  t := time.Now()
	n := "web"
	a := true

	h := github.Hook{
		CreatedAt:  &t,
		UpdatedAt: &t,
		Name:      &n,
		Events:    []string{"push"},
		Config:     map[string]interface{}{
			"url":     r.Trigger(),
			"content_type": "json",
		},
		Active: &a,
	}

	repoName, err := r.GetShortName()

	if err != nil {
		return "", err
	}

	hk, response, err := client.Repositories.CreateHook(r.GetUserName(),repoName, &h)
	if err != nil {
		return "", err
	}

  m.debugResp(response.Response)

	log.Debugf("[github] created webhook [%s,%s] successfully.",  r.Gitr(), strconv.Itoa(*hk.ID))
	//We need to save the hook id.
	return strconv.Itoa(*hk.ID), nil

}

func (m githubManager) RemoveHook(r repository.Repository) error {
	//get  a single hook saved id and remove the hook.
	return nil
}

func (m githubManager) debugResp(resp *http.Response) {
  log.Debugf(cmd.Colorfy("--- git --", "yellow", "",""))
	responseDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return
	}

	log.Debugf("%v",responseDump)
}
