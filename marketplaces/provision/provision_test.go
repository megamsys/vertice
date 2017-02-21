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
package provision

/*
import (
	"errors"
	"reflect"
	"testing"

	"gopkg.in/check.v1"
)

type ProvisionSuite struct{}

var _ = check.Suite(ProvisionSuite{})

func Test(t *testing.T) {
	check.TestingT(t)
}

func (ProvisionSuite) TestRegisterAndGetProvisioner(c *check.C) {
	var p Provisioner
	Register("my-provisioner", p)
	got, err := Get("my-provisioner")
	c.Assert(err, check.IsNil)
	c.Check(got, check.DeepEquals, p)
	_, err = Get("unknown-provisioner")
	c.Check(err, check.NotNil)
	expectedMessage := `unknown provisioner: "unknown-provisioner"`
	c.Assert(err.Error(), check.Equals, expectedMessage)
}

func (ProvisionSuite) TestRegistry(c *check.C) {
	var p1, p2 Provisioner
	Register("my-provisioner", p1)
	Register("your-provisioner", p2)
	provisioners := Registry()
	alt1 := []Provisioner{p1, p2}
	alt2 := []Provisioner{p2, p1}
	if !reflect.DeepEqual(provisioners, alt1) && !reflect.DeepEqual(provisioners, alt2) {
		c.Errorf("Registry(): Expected %#v. Got %#v.", alt1, provisioners)
	}
}

func (ProvisionSuite) TestError(c *check.C) {
	errs := []*Error{
		{Reason: "something", Err: errors.New("went wrong")},
		{Reason: "something went wrong"},
	}
	expected := []string{"went wrong: something", "something went wrong"}
	for i := range errs {
		c.Check(errs[i].Error(), check.Equals, expected[i])
	}
}

func (ProvisionSuite) TestErrorImplementsError(c *check.C) {
	var _ error = &Error{}
}


func (ProvisionSuite) TestStatusString(c *check.C) {
	var s Status = "pending"
	c.Assert(s.String(), check.Equals, "pending")
}

func (ProvisionSuite) TestStatuses(c *check.C) {
	c.Check(StatusBootstrapped.String(), check.Equals, "bootstrapped")
	c.Check(StatusBuilding.String(), check.Equals, "building")
	c.Check(StatusError.String(), check.Equals, "error")
	c.Check(StatusStarted.String(), check.Equals, "started")
	c.Check(StatusStopped.String(), check.Equals, "stopped")
	c.Check(StatusStarting.String(), check.Equals, "starting")
}


func (ProvisionSuite) TestBoxAvailable(c *check.C) {
	var tests = []struct {
		input    Status
		expected bool
	}{
		{StatusBootstrapped, false},
		{StatusStarting, true},
		{StatusStarted, true},
		{StatusBuilding, false},
		{StatusError, true},
	}
	for _, test := range tests {
		b := Box{Status: test.input}
		c.Check(b.Available(), check.Equals, test.expected)
	}
}
*/
