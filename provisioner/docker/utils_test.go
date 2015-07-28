package docker

import (
	"net"
	"gopkg.in/check.v1"
	"testing"
)

type S struct{}

func Test(t *testing.T) {
	check.TestingT(t)
}


func (s *S) TestGetIpFullMask(c *check.C) {

	_, subnet, _ := net.ParseCIDR("192.168.1.89/24")
	_, _, err := IPRequest(*subnet)
	
	c.Assert(err, check.IsNil)
}

