package deployd

import (
	"github.com/megamsys/megamd/meta"
	"github.com/megamsys/megamd/services/deployd"
	"gopkg.in/check.v1"
)

type S struct {
	service *deployd.Service
}

var _ = check.Suite(&S{})

func (s *S) SetUpSuite(c *check.C) {
	srv, err := &NewService(deployd.Config{
		BindAddress: "127.0.0.1:0",
	})
	s.service = srv
	c.Assert(err, check.IsNil)
}
