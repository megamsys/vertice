package route53

import (
	"fmt"
	"os"
	"testing"

	"github.com/megamsys/megamd/router"
	"github.com/megamsys/megamd/subd/dns"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct {
	cf *dns.Config
}

var _ = check.Suite(&S{})

var aws_acc = os.Getenv("AWS_ACCESS_KEY")
var aws_sec = os.Getenv("AWS_SECRET_KEY")

func (s *S) SetUpSuite(c *check.C) {
	if aws_acc == "" || aws_sec == "" {
		c.Skip("-R53 (aws access/secret keys) not provided")
	}
	cf := dns.NewConfig()
	cf.AccessKey = os.Getenv("AWS_ACCESS_KEY")
	cf.SecretKey = os.Getenv("AWS_SECRET_KEY")
	s.cf = cf
}

func (s *S) TestShouldBeRegistered(c *check.C) {
	s.cf.MkGlobal()
	router.Register("route53", createRouter)
	got, err := router.Get("route53")
	c.Assert(err, check.IsNil)
	_, ok := got.(route53Router)
	c.Assert(ok, check.Equals, true)
}

func (s *S) TestSetCName(c *check.C) {
	s.cf.MkGlobal()
	vRouter, err := router.Get("route53")
	c.Assert(err, check.IsNil)
	err = vRouter.SetCName("myapp1.megambox.com", "192.168.1.100")
	c.Assert(err, check.IsNil)
}

func (s *S) TestAddr(c *check.C) {
	s.cf.MkGlobal()
	vRouter, err := router.Get("route53")
	c.Assert(err, check.IsNil)
	addr, err := vRouter.Addr("myapp1.megambox.com")
	c.Assert(err, check.IsNil)
	c.Assert(addr, check.Equals, "megambox.com.")
}
func (s *S) TestAddrNotExist(c *check.C) {
	s.cf.MkGlobal()
	vRouter, err := router.Get("route53")
	c.Assert(err, check.IsNil)
	addr, err := vRouter.Addr("myapp.megamboxy.com")
	c.Assert(addr, check.Equals, "")
}

func (s *S) TestSetCNameDuplicate(c *check.C) {
	s.cf.MkGlobal()
	vRouter, err := router.Get("route53")
	err = vRouter.SetCName("myapp1.megambox.com", "192.168.1.100")
	c.Assert(err, check.Equals, nil)
}

func (s *S) TestUnsetCName(c *check.C) {
	s.cf.MkGlobal()
	vRouter, err := router.Get("route53")
	c.Assert(err, check.IsNil)
	err = vRouter.UnsetCName("myapp1.megambox.com", "192.168.1.100")
	c.Assert(err, check.IsNil)
}

func (s *S) TestUnsetCNameNotExist(c *check.C) {
	s.cf.MkGlobal()
	vRouter, err := router.Get("route53")
	c.Assert(err, check.IsNil)
	err = vRouter.UnsetCName("myapp2.megambox66.com", "192.168.1.102")
	c.Assert(err, check.Equals, router.ErrDomainNotFound)
}

func (s *S) TestStartupMessage(c *check.C) {
	got, err := router.Get("route53")
	c.Assert(err, check.IsNil)
	mRouter, ok := got.(route53Router)
	c.Assert(ok, check.Equals, true)
	message, err := mRouter.StartupMessage()
	c.Assert(err, check.IsNil)
	c.Assert(message, check.Equals, fmt.Sprintf("R53 router ok!"))
}
