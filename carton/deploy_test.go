/*
** copyright [2013-2016] [Megam Systems]
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

/*
import (
	"gopkg.in/check.v1"
)

func (s *S) TestDeployToProvisioner(c *check.C) {
	carton := provisiontest.NewFakeCarton("myapp", "tosca.torpedo.ubuntu", provision.BoxNone, 1)
	defer s.provisioner.Destroy(&box)
	writer := &bytes.Buffer{}
	opts := DeployOpts{B: &carton.Boxs()[0]}
	_, err = deployToProvisioner(&opts, writer)
	c.Assert(err, check.IsNil)
	logs := writer.String()
	c.Assert(logs, check.Equals, "Git deploy called")
}

func (s *S) TestDeployToProvisionerImage(c *check.C) {
	carton := provisiontest.NewFakeCarton("myapp", "tosca.web.java", provision.BoxSome, 1)
	defer s.provisioner.Destroy(&box)
	writer := &bytes.Buffer{}
	opts := DeployOpts{B: &carton.Boxs()[0]}
	_, err = deployToProvisioner(&opts, writer)
	c.Assert(err, check.IsNil)
	logs := writer.String()
	c.Assert(logs, check.Equals, "Image deploy called")
}
*/
