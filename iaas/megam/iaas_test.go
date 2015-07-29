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
package megam

import (
	"gopkg.in/check.v1"
	"github.com/megamsys/megamd/global"
	"testing"
	"encoding/json"
	"github.com/tsuru/config"
)

type S struct{}

const defaultConfigPath = "conf/megamd.conf"

var _ = check.Suite(&S{})

var megam = MegamIaaS{}	
	
var json1 = `{"key":"domain","value":"megambox.com"}`

	var res1 = &global.KeyValuePair{}
	var input1 = json.Unmarshal([]byte(json1), &res1)
	
	var json2 = `{"key":"sshkey","value":"ssss"}`

	var res2 = &global.KeyValuePair{}
	var input2 = json.Unmarshal([]byte(json2), &res2)
	
	var json3 = `{"key":"provider","value":"chef"}`

	var res3 = &global.KeyValuePair{}
	var input3 = json.Unmarshal([]byte(json3), &res3)
	
	var json4 = `{"key":"version","value":"8"}`

	var res4 = &global.KeyValuePair{}
	var input4 = json.Unmarshal([]byte(json4), &res4)
	
	var json5 = `{"key":"cpu","value":"8"}`

	var res5 = &global.KeyValuePair{}
	var input5 = json.Unmarshal([]byte(json5), &res5)
	
	var json6 = `{"key":"ram","value":"8"}`

	var res6 = &global.KeyValuePair{}
	var input6 = json.Unmarshal([]byte(json6), &res6)
	
	var inputs = make([]*global.KeyValuePair, 0)
	var d = append(inputs, res1)
	var e = append(d, res2)
	var f = append(e, res3)
	var g = append(f, res4)
	var h = append(g, res5)
	var i = append(h, res6)
		
var	asm = &global.AssemblyWithComponents{
							Id: 			"ASM1226483770383794176",
							Name: 			"winfred", 
							ToscaType: 		"tosca.torpedo.debian",
							Inputs:			i,
							Status:			"LAUNCHING",
							CreatedAt:		"2015-06-27 07:38:51 +0000",
							}

func Test(t *testing.T) {
	check.TestingT(t)
}

func (s *S) TestCreateMachine(c *check.C) {
	setconfig()
	defer unsetconfig()
			
	_, err := megam.CreateMachine(&global.PredefClouds{}, asm, "ACT1226153924164190208")
	c.Assert(err, check.IsNil)
}

/*
* set the config variables for testing
*/
func setconfig() {
	config.Set("opennebula:access_key", "accesskey")
	config.Set("opennebula:secret_key", "secret_key")
	config.Set("opennebula:ssh_user", "ssh_user")
	config.Set("opennebula:identity_file", "identity_file")
	config.Set("opennebula:zone", "zone")
	config.Set("knife:path", "/var/lib/megam/megamd/chef-repo/.chef/knife.rb")
	config.Set("knife:recipe", "megam_run")
	config.Set("launch:riak", "api.megam.io")
	config.Set("launch:rabbitmq", "api.megam.io")
	config.Set("launch:monitor", "api.megam.io")
	config.Set("launch:kibana", "api.megam.io")
	config.Set("launch:etcd", "api.megam.io")
}

/*
* unset the config variables for testing
*/
func unsetconfig() {
	config.Unset("opennebula:access_key")	
	config.Unset("opennebula:secret_key")
	config.Unset("opennebula:ssh_user")
	config.Unset("opennebula:identity_file")
	config.Unset("opennebula:zone")
	config.Unset("knife:path")
	config.Unset("knife:recipe")
	config.Unset("launch:riak")
	config.Unset("launch:rabbitmq")
	config.Unset("launch:monitor")
	config.Unset("launch:kibana")
	config.Unset("launch:etcd")
}

