package carton


import (
	"github.com/megamsys/libgo/api"
	"testing"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct {
	Credentials api.ApiArgs
	Master_Key string
	Host  string
	//provisioner *onetest.FakeOneProvisioner
}

var _ = check.Suite(&S{})

//we need make sure the stub deploy methods are supported.
func (s *S) SetUpSuite(c *check.C) {
	s.Credentials = api.ApiArgs{
		Master_Key: "3b8eb672aa7c8db82e5d34a0744740b20ed59e1f6814cfb63364040b0994ee3f",
		Url: "http://40.74.121.55:9000/v2",
	}
}
