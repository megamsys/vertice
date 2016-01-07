package machine

import (
	/*	"bytes"

		"github.com/megamsys/megamd/provision"
		"github.com/megamsys/megamd/provision/provisiontest"
		"github.com/megamsys/opennebula-go/compute" */
	"gopkg.in/check.v1"
)

func (s *S) TestMachineName(c *check.C) {
	mach := Machine{Id: "abc123", Name: "alpha.megambox.com"}
	c.Check(mach.Id, check.Equals, "abc123")
	c.Check(mach.Name, check.Equals, "alpha.megambox.com")
}

/*
this needs OneServer fix.
func (s *S) TestMachineCreate(c *check.C) {
	carton := provisiontest.NewFakeCarton("abdulkalam.megambox.com", "tosca.torpedo.ubuntu", provision.BoxSome, 0)

	mach := Machine{
		Name:     "abdulkalam.megambox.com",
		Id:       "CMP010101010101",
		CartonId: "ASM010101010101",
		Level:    provision.BoxSome,
		Image:    "ubuntu",
		Routable: false,
		Status:   provision.StatusLaunching,
	}

	err := mach.Create(&CreateArgs{
		Box:         nil,
		Compute:     carton.Compute,
		Deploy:      true,
		Provisioner: s.p,
	})
	c.Assert(err, check.IsNil)

}

func (s *S) TestMachineRemove(c *check.C) {
	_ = provisiontest.NewFakeCarton("abdulkalam.megambox.com", "tosca.torpedo.ubuntu", provision.BoxSome, 0)

	mach := Machine{
		Name:     "abdulkalam.megambox.com",
		Id:       "CMP010101010101",
		CartonId: "ASM010101010101",
		Level:    provision.BoxSome,
		Image:    "ubuntu",
		Routable: false,
		Status:   provision.StatusLaunching,
	}
	err := mach.Remove(s.p)
	c.Assert(err, check.IsNil)
}

//this fails, needs OneServer to return back nil
func (s *S) TestMachineLogs(c *check.C) {
	opts := compute.VirtualMachine{}
	_, err := s.newMachine(opts, nil)
	c.Assert(err, check.NotNil)
	c.Assert(mach,check.NotNil)
		var buff bytes.Buffer
		err = mach.Logs(s.p, &buff)
		c.Assert(err, check.IsNil)
		c.Assert(buff.String(), check.Not(check.Equals), "")
}
*/
