package marketplaces

import (
	"github.com/megamsys/vertice/meta"
	"testing"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct {
}

var _ = check.Suite(&S{})

//we need make sure the stub deploy methods are supported.
func (s *S) SetUpSuite(c *check.C) {
	mc := meta.Config{
		Api:  "http://192.168.0.14:9000/v2",
		MasterUser: "info@megam.io",
		MasterKey: "3b8eb672aa7c8db82e5d34a0744740b20ed59e1f6814cfb63364040b0994ee3f",
	}
  mc.MkGlobal()
}
