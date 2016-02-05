package docker

import ()

/*
func (s *S) TestUsePlatformImage(c *check.C) {
	app1 := &app.App{Name: "app1", Platform: "python", Deploys: 40}
	err := s.storage.Apps().Insert(app1)
	c.Assert(err, check.IsNil)
	ok := s.p.usePlatformImage(app1)
	c.Assert(ok, check.Equals, true)
	app2 := &app.App{Name: "app2", Platform: "python", Deploys: 20}
	err = s.storage.Apps().Insert(app2)
	c.Assert(err, check.IsNil)
	ok = s.p.usePlatformImage(app2)
	c.Assert(ok, check.Equals, true)
	app3 := &app.App{Name: "app3", Platform: "python", Deploys: 0}
	err = s.storage.Apps().Insert(app3)
	c.Assert(err, check.IsNil)
	ok = s.p.usePlatformImage(app3)
	c.Assert(ok, check.Equals, true)
	app4 := &app.App{Name: "app4", Platform: "python", Deploys: 19}
	err = s.storage.Apps().Insert(app4)
	c.Assert(err, check.IsNil)
	ok = s.p.usePlatformImage(app4)
	c.Assert(ok, check.Equals, false)
	app5 := &app.App{
		Name:           "app5",
		Platform:       "python",
		Deploys:        19,
		UpdatePlatform: true,
	}
	err = s.storage.Apps().Insert(app5)
	c.Assert(err, check.IsNil)
	ok = s.p.usePlatformImage(app5)
	c.Assert(ok, check.Equals, true)
	app6 := &app.App{Name: "app6", Platform: "python", Deploys: 19}
	err = s.storage.Apps().Insert(app6)
	c.Assert(err, check.IsNil)
	coll := s.p.Collection()
	defer coll.Close()
	err = coll.Insert(container.Container{AppName: app6.Name, Image: "github.com/megamsys/vertice/app-app6"})
	c.Assert(err, check.IsNil)
	ok = s.p.usePlatformImage(app6)
	c.Assert(ok, check.Equals, false)
}


func (s *S) TestPlatformImageName(c *check.C) {
	platName := platformImageName("python")
	c.Assert(platName, check.Equals, "github.com/megamsys/vertice/python:latest")
	config.Set("docker:registry", "localhost:3030")
	defer config.Unset("docker:registry")
	platName = platformImageName("ruby")
	c.Assert(platName, check.Equals, "localhost:3030/"github.com/megamsys/vertice/ruby:latest")
}
*/
