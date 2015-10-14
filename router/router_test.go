
package router

import (
	"gopkg.in/check.v1"
)


func (s *S) TestRegisterAndGet(c *check.C) {
	var r Router
	var prefixes []string
	routerCreator := func(prefix string) (Router, error) {
		prefixes = append(prefixes, prefix)
		return r, nil
	}
	Register("mine", routerCreator)
	got, err := Get("mine")
	c.Assert(err, check.IsNil)
	c.Assert(r, check.DeepEquals, got)
}

func (s *S) TestRegisterAndGetCustomNamedRouter(c *check.C) {
	var prefixes []string
	routerCreator := func(prefix string) (Router, error) {
		prefixes = append(prefixes, prefix)
		var r Router
		return r, nil
	}
	Register("myrouter", routerCreator)
	_, err := Get("myrouter")
	c.Assert(err, check.IsNil)	
}
