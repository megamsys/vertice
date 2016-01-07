package machine

import (
	"testing"

	"github.com/megamsys/megamd/provision"
	"github.com/megamsys/opennebula-go/compute"
	otesting "github.com/megamsys/opennebula-go/testing"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

var _ = check.Suite(&S{})

type S struct {
	p      *fakeOneProvisioner
	server *otesting.OneServer
}

func (s *S) SetUpSuite(c *check.C) {
}

/*func (s *S) SetUpTest(c *check.C) {
	server, err := otesting.NewServer("127.0.0.1:5555")
	c.Assert(err, check.IsNil)
	s.server = server
	s.p, err = newFakeOneProvisioner(s.server.URL())
	c.Assert(err, check.IsNil)
}

func (s *S) TestTearDownTest(c *check.C) {
	s.server.Stop()
}*/

func (s *S) newMachine(opts compute.VirtualMachine, p *fakeOneProvisioner) (*Machine, error) {
	if p == nil {
		p = s.p
	}

	mach := Machine{
		Name:     "screwapp.megambox.com",
		Id:       "CMP010101010101",
		CartonId: "ASM010101010101",
		Level:    provision.BoxSome,
		Image:    "figureoutimage",
		Routable: false,
		Status:   provision.StatusLaunching,
	}

	_, _, err := p.Cluster().CreateVM(opts)
	if err != nil {
		return nil, err
	}
	return &mach, nil
}
