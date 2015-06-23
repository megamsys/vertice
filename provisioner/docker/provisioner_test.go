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

package docker

import (
	"encoding/json"
	"testing"

	"github.com/megamsys/megamd/global"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type S struct{}

var _ = check.Suite(&S{})

var jsonC = `{"id": "COM1225063232952205312","name":"canards","tosca_type":"tosca.addon.containers","inputs":[{"key":"domain","value":"megambox.com"},{"key":"version","value":"0.1.0"},{"key":"source","value":"notgood/container"}],"outputs":[],"artifacts":{"artifact_type":"","content":"","artifact_requirements":[]},"related_components":[],"operations":[],"status":"LAUNCHING","created_at":"2015-06-23 09:34:09 +0000"}`
var res = &global.Component{}
var components = json.Unmarshal([]byte(jsonC), &res)

var jsonI = `{"key":"endpoint","value":"tcp://103.56.92.2:2375"}`

var resI = &global.KeyValuePair{}
var input1 = json.Unmarshal([]byte(jsonI), &resI)

var comp = make([]*global.Component, 0)
var a = append(comp, res)

var inputs = make([]*global.KeyValuePair, 0)
var b = append(inputs, resI)

func (s *S) TestDockerCreate(c *check.C) {
	asm := &global.AssemblyWithComponents{Id: "ASM1225063232943816704", Name: "spindlier", ToscaType: "tosca.addon.ubuntu",
		Components: a,
		Inputs:     b}

	dockertest := &Docker{}

	_, err := dockertest.Create(asm, "RIP1225063233107394560", false, "ACT1225004185968312320")
	c.Assert(err, check.IsNil)

}

func (s *S) TestFailDockerCreate(c *check.C) {

	var jsonI = `{"key":"endpoint","value":"tcp://103.56333.92.2:2375"}`

	var resI = &global.KeyValuePair{}
	json.Unmarshal([]byte(jsonI), &resI)
	var inputs = make([]*global.KeyValuePair, 0)
	var b = append(inputs, resI)

	asm := &global.AssemblyWithComponents{Id: "ASM1225063232943816704", Name: "spindlier", ToscaType: "tosca.addon.ubuntu",
		Components: a,
		Inputs:     b}

	dockertest := &Docker{}

	_, err := dockertest.Create(asm, "RIP1225063233107394560", false, "ACT1225004185968312320")
	c.Assert(err, check.IsNil)

}
