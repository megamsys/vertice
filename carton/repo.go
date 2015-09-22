package carton

import (
	"strconv"

	"github.com/megamsys/megamd/repository"
)

func NewRepo(ops []*Operations, optype string) repository.Repo {
	o := parseOps(ops, optype)

	if o != nil {
		enabled, _ := strconv.ParseBool(o.OperationRequirements.match(repository.CI_ENABLED))

		return repository.Repo{
			Enabled:  enabled,
			Type:     o.OperationRequirements.match(repository.CI_TYPE),
			Token:    o.OperationRequirements.match(repository.CI_TOKEN),
			Source:   o.OperationRequirements.match(repository.CI_SOURCE),
			UserName: o.OperationRequirements.match(repository.CI_USER),
			GitURL:   o.OperationRequirements.match(repository.CI_URL),
		}

	}
	return repository.Repo{}
}

func parseOps(ops []*Operations, optype string) *Operations {
	for _, o := range ops {
		switch o.OperationType {
		case repository.CI:
			return o
		}
	}
	return nil
}
