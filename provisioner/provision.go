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
package provisioner

import (
	"fmt"

	"github.com/megamsys/megamd/global"
)

/*
 * Provisioner is the basic interface of this package.
 */
type Provisioner interface {
	Create(*global.AssemblyWithComponents, string, bool, string) (string, error)
	Delete(*global.AssemblyWithComponents, string) (string, error)
}

var provisioners = make(map[string]Provisioner)

/*
 * Register registers a new provisioner in the Provisioner registry.
 */

func Register(name string, p Provisioner) {
	provisioners[name] = p
}

func GetProvisioner(name string) (Provisioner, error) {

	provider, ok := provisioners[name]
	if !ok {
		return nil, fmt.Errorf("Provisioner not registered")
	}
	return provider, nil
}
