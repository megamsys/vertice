/*
** Copyright [2013-2016] [Megam Systems]
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
package carton

import (
	"github.com/megamsys/megamd/repository"
)

type Operations struct {
	Type        string    `json:"operation_type"`
	Description string    `json:"description"`
	Properties  JsonPairs `json:"properties"`
	Status      string    `json:"status"`
}

func (o *Operations) prepBuildHook() *repository.Hook {
	return &repository.Hook{
		Enabled:  true,
		Token:    o.Properties.match(repository.TOKEN),
		UserName: o.Properties.match(repository.USERNAME),
	}
}

func BuildHook(ops []*Operations, opsType string) *repository.Hook {
	for _, o := range ops {
		switch o.Type {
		case opsType:
			return o.prepBuildHook()
		}
	}
	return nil
}
