package route53
/*
import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/megamsys/megamd/router"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type S struct {
	c *dns.Config
}

var _ = check.Suite(&S{})

func (s *S) SetUpSuite(c *check.C) {
	c := dns.NewConfig()
	if _, err := toml.Decode(`
enabled = true
access_key  = "temp_access_key"
secrete_key = "temp_secret_key"
`, &d); err != nil {
		c.Fatal(err)
	}
	s.c = c
}


func (s *S) TestShouldBeRegistered(c *check.C) {
	got, err := router.Get("vulcand")
	c.Assert(err, check.IsNil)
	r, ok := got.(*vulcandRouter)
	c.Assert(ok, check.Equals, true)
	c.Assert(r.client.Addr, check.Equals, s.vulcandServer.URL)
	c.Assert(r.domain, check.Equals, "vulcand.example.com")
}

func (s *S) TestSetCName(c *check.C) {
	vRouter, err := router.Get("vulcand")
	c.Assert(err, check.IsNil)
	err = vRouter.AddBackend("myapp")
	c.Assert(err, check.IsNil)
	u1, _ := url.Parse("http://1.1.1.1:111")
	u2, _ := url.Parse("http://2.2.2.2:222")
	err = vRouter.AddRoute("myapp", u1)
	c.Assert(err, check.IsNil)
	err = vRouter.AddRoute("myapp", u2)
	c.Assert(err, check.IsNil)
	err = vRouter.SetCName("myapp.cname.example.com", "myapp")
	c.Assert(err, check.IsNil)
	appFrontend, err := s.engine.GetFrontend(engine.FrontendKey{
		Id: "tsuru_myapp.vulcand.example.com",
	})
	c.Assert(err, check.IsNil)
	cnameFrontend, err := s.engine.GetFrontend(engine.FrontendKey{
		Id: "tsuru_myapp.cname.example.com",
	})
	c.Assert(err, check.IsNil)
	c.Assert(cnameFrontend.BackendId, check.DeepEquals, appFrontend.BackendId)
	c.Assert(cnameFrontend.Route, check.Equals, `Host("myapp.cname.example.com")`)
	c.Assert(cnameFrontend.Type, check.Equals, "http")
}

func (s *S) TestSetCNameDuplicate(c *check.C) {
	vRouter, err := router.Get("vulcand")
	c.Assert(err, check.IsNil)
	err = vRouter.AddBackend("myapp")
	c.Assert(err, check.IsNil)
	u1, _ := url.Parse("http://1.1.1.1:111")
	u2, _ := url.Parse("http://2.2.2.2:222")
	err = vRouter.AddRoute("myapp", u1)
	c.Assert(err, check.IsNil)
	err = vRouter.AddRoute("myapp", u2)
	c.Assert(err, check.IsNil)
	err = vRouter.SetCName("myapp.cname.example.com", "myapp")
	c.Assert(err, check.IsNil)
	err = vRouter.SetCName("myapp.cname.example.com", "myapp")
	c.Assert(err, check.Equals, router.ErrCNameExists)
}

func (s *S) TestUnsetCName(c *check.C) {
	vRouter, err := router.Get("vulcand")
	c.Assert(err, check.IsNil)
	err = vRouter.AddBackend("myapp")
	c.Assert(err, check.IsNil)
	u1, _ := url.Parse("http://1.1.1.1:111")
	u2, _ := url.Parse("http://2.2.2.2:222")
	err = vRouter.AddRoute("myapp", u1)
	c.Assert(err, check.IsNil)
	err = vRouter.AddRoute("myapp", u2)
	c.Assert(err, check.IsNil)
	err = vRouter.SetCName("myapp.cname.example.com", "myapp")
	c.Assert(err, check.IsNil)
	frontends, err := s.engine.GetFrontends()
	c.Assert(err, check.IsNil)
	c.Assert(frontends, check.HasLen, 2)
	vRouter.UnsetCName("myapp.cname.example.com", "myapp")
	frontends, err = s.engine.GetFrontends()
	c.Assert(err, check.IsNil)
	c.Assert(frontends, check.HasLen, 1)
	c.Assert(frontends[0].Id, check.Equals, "tsuru_myapp.vulcand.example.com")
}

func (s *S) TestUnsetCNameNotExist(c *check.C) {
	vRouter, err := router.Get("vulcand")
	c.Assert(err, check.IsNil)
	frontends, err := s.engine.GetFrontends()
	c.Assert(err, check.IsNil)
	c.Assert(frontends, check.HasLen, 0)
	err = vRouter.UnsetCName("myapp.cname.example.com", "myapp")
	c.Assert(err, check.Equals, router.ErrCNameNotFound)
}

func (s *S) TestAddr(c *check.C) {
	vRouter, err := router.Get("vulcand")
	c.Assert(err, check.IsNil)
	err = vRouter.AddBackend("myapp")
	c.Assert(err, check.IsNil)
	addr, err := vRouter.Addr("myapp")
	c.Assert(err, check.IsNil)
	c.Assert(addr, check.Equals, "myapp.vulcand.example.com")
}

func (s *S) TestAddrNotExist(c *check.C) {
	vRouter, err := router.Get("vulcand")
	c.Assert(err, check.IsNil)
	frontends, err := s.engine.GetFrontends()
	c.Assert(err, check.IsNil)
	c.Assert(frontends, check.HasLen, 0)
	backends, err := s.engine.GetBackends()
	c.Assert(err, check.IsNil)
	c.Assert(backends, check.HasLen, 0)
	addr, err := vRouter.Addr("myapp")
	c.Assert(err, check.Equals, router.ErrBackendNotFound)
	c.Assert(addr, check.Equals, "")
}

func (s *S) TestStartupMessage(c *check.C) {
	got, err := router.Get("vulcand")
	c.Assert(err, check.IsNil)
	mRouter, ok := got.(router.MessageRouter)
	c.Assert(ok, check.Equals, true)
	message, err := mRouter.StartupMessage()
	c.Assert(err, check.IsNil)
	c.Assert(message, check.Equals,
		fmt.Sprintf(`vulcand router "vulcand.example.com" with API at "%s"`, s.vulcandServer.URL),
	)
}
*/
