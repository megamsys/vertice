package machine

import (
//	"github.com/megamsys/megamd/provision"
//	"github.com/megamsys/megamd/provision/provisiontest"
	"gopkg.in/check.v1"
)

func (s *S) TestMachineName(c *check.C) {
	mach := Machine{Id: "abc123", Name: "alpha.megambox.com"}
	c.Check(mach.Id, check.Equals, "abc123")
	c.Check(mach.Name, check.Equals, "alpha.megambox.com")
}
/*

func (s *S) TestMachineCreate(c *check.C) {
	carton := provisiontest.NewFakeCarton("dabba.megambox.com", "tosca.torpedo.ubuntu", provision.BoxSome, 0)

	mach := Machine{
		Name:     "kalam.megambox.com",
		Id:       "CMP010101010101",
		CartonId: "ASM010101010101",
		Level:    provision.BoxSome,
		Image:    "ubuntu",
		Routable: false,
		Status:   provision.StatusDeploying,
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
	app := provisiontest.NewFakeCarton("app-name", "", 1)
	app.Memory = 15
	app.Swap = 15
	app.CpuShare = 50
	mach := Machine{
		Name:        "myName",
		AppName:     app.GetName(),
		Type:        app.GetPlatform(),
		Status:      "created",
		ProcessName: "myprocess1",
	}
	err := mach.Remove(s.p)
	c.Assert(err, check.IsNil)
}

func (s *S) TestMachineLogs(c *check.C) {
	mach, err := s.newMachine(newMachineOpts{}, nil)
	c.Assert(err, check.IsNil)
	var buff bytes.Buffer
	err = mach.Logs(s.p, &buff)
	c.Assert(err, check.IsNil)
	c.Assert(buff.String(), check.Not(check.Equals), "")
}
*/
