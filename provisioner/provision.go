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
	//	"errors"
	"fmt"
	"github.com/megamsys/megamd/global"
)

// Provisioner is the basic interface of this package.
//
type Provisioner interface {
	// Provision is called when tsuru is creating the app.
	//Provision(*global.AssemblyWithComponents) error

	// ExecuteCommand runs a command in all units of the app.
	 CreateCommand(*global.AssemblyWithComponents, string, bool, string) (string, error)
     DeleteCommand(*global.AssemblyWithComponents, string) (string, error)
	// ExecuteCommandOnce runs a command in one unit of the app.
	//	ExecuteCommandOnce(stdout, stderr io.Writer, app global.AssemblyWithComponents, cmd string, args ...string) error

}

var provisioners = make(map[string]Provisioner)

// Register registers a new provisioner in the Provisioner registry.
func Register(name string, p Provisioner) {
	provisioners[name] = p
}

func GetProvisioner(name string) (Provisioner, error) {

	provider, ok := provisioners[name]
	if !ok {
		return nil, fmt.Errorf("Provisioner not registered")
	}
	return provider, nil
	//return nil
}
