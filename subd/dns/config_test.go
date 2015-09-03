package dns

import (
	"github.com/BurntSushi/toml"
	"github.com/megamsys/megamd/services/dns"
	"gopkg.in/check.v1"
)

// Ensure the configuration can be parsed.
func (s *S) TestConfig_Parse(c *check.C) {
	// Parse configuration.
	var c dns.Config
	if _, err := toml.Decode(`
enabled = true
route53_access_key  = ":9000"
route53_secrete_key = "xxx"
`, &c); err != nil {
		t.Fatal(err)
	}

	c.Assert(c.Route53SecretKey, check.Equals, "locahost")
	c.Assert(c.Route53AccessKey, check.Equals, "locahost")
}
