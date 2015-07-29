package gitlab

import (
	log "code.google.com/p/log4go"
	"github.com/megamsys/megamd/global"
	"github.com/megamsys/megamd/plugins"
	"github.com/plouc/go-gitlab-client"
	"github.com/tsuru/config"
)

/*
 * Init to register a plugin
 */

func Init() {
	plugins.RegisterPlugins("gitlab", &GitlabPlugin{})
}

type GitlabPlugin struct{}

const (
	GITLAB = "gitlab"
	ENABLE = "true"
)

func (c *GitlabPlugin) Watcher(asm *global.AssemblyWithComponents, ci *global.Operations, com *global.Component) error {
	switch ci.OperationType {
	case "CI":
		cierr := cioperation(asm, ci, com)
		if cierr != nil {
			return cierr
		}
		break
	}
	return nil
}

/* GITLAB CE Support - Gitlab - Private git repository
* cioperation builds the TRIGGERURL and calls gitlab client to add a webhook.
 */
func cioperation(asm *global.AssemblyWithComponents, ci *global.Operations, com *global.Component) error {

	pair_scm, perrscm := global.ParseKeyValuePair(ci.OperationRequirements, "ci-scm")
	if perrscm != nil {
		log.Error("Failed to get the domain value : %s", perrscm)
	}

	pair_enable, perrenable := global.ParseKeyValuePair(ci.OperationRequirements, "ci-enable")
	if perrenable != nil {
		log.Error("Failed to get the domain value : %s", perrenable)
	}

	if pair_scm.Value == GITLAB && pair_enable.Value == ENABLE {
		log.Info("GitLab is working..")

		pair_token, perrtoken := global.ParseKeyValuePair(ci.OperationRequirements, "ci-token")
		if perrtoken != nil {
			log.Error("Failed to get the ci-token value : %s", perrtoken)

		}

		pair_url, perrtoken := global.ParseKeyValuePair(ci.OperationRequirements, "ci-url")
		if perrtoken != nil {
			log.Error("Failed to get the ci-url value : %s", perrtoken)

		}
		pair_apiversion, perrtoken := global.ParseKeyValuePair(ci.OperationRequirements, "ci-apiversion")
		if perrtoken != nil {
			log.Error("Failed to get the ci-apiversion value : %s", perrtoken)

		}

		pair_owner, perrowner := global.ParseKeyValuePair(ci.OperationRequirements, "ci-owner")
		if perrowner != nil {
			log.Error("Failed to get the ci-owner value : %s", perrowner)
		}

		api_host, apierr := config.GetString("megam:api")
		if apierr != nil {
			return apierr
		}

		trigger_url :=  api_host +  "/assembly/build/" + asm.Id + "/" + com.Id

		client := gogitlab.NewGitlab(pair_url.Value, pair_apiversion.Value, pair_token.Value)

		err := client.AddProjectHook(pair_owner.Value, trigger_url, false, false, false)
		if err != nil {
			return err
		}
		log.Info("[megamd] added project hook %s %s.",pair_owner.Value, trigger_url)
	} else {
		log.Info("[megamd] skip gitlab.")
	}
	return nil

}

func (c *GitlabPlugin) Notify(m *global.EventMessage) error { return nil }
