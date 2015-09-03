package docker

import (
	"github.com/megamsys/megamd/services/docker"
	"gopkg.in/check.v1"

)

type S struct {
	service *docker.Service
}

var _ = check.Suite(&S{})

func (s *S) SetUpSuite(c *check.C) {
	srv, err := &NewService(docker.Config{
		BindAddress: "127.0.0.1:0",
	})
	s.service = srv
	c.Assert(err, check.IsNil)
}
