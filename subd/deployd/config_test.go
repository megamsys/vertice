package deployd

import (
	"github.com/BurntSushi/toml"
	"gopkg.in/check.v1"
)

// Ensure the configuration can be parsed.
func (s *S) TestConfig_Parse(c *check.C) {
	// Parse configuration.
	var c deployd.Config
	if _, err := toml.Decode(`
		one_endpoint = "http://opennebula:3000/xmlrpc2"
		one_userid   = "oneadmin"
		one_password = "password"
		one_template = "megam_trusty"
		one_zone     = "plano01"
		certificate = "/etc/ssl/cert.pem"

`, &c); err != nil {
		t.Fatal(err)
	}

	c.Assert(c.OneEndPoint, check.Equals, "locahost")
	c.Assert(c.OneUserid, check.Equals, "locahost")
	c.Assert(c.OnePassword, check.Equals, "locahost")
}
