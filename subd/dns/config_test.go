package dns

import (
	"github.com/BurntSushi/toml"
	"gopkg.in/check.v1"
)

type S struct {
	//	service *deployd.Service
}

var _ = check.Suite(&S{})

// Ensure the configuration can be parsed.
func (s *S) TestDnsConfig_Parse(c *check.C) {
	// Parse configuration.
	var cm Config
	if _, err := toml.Decode(`
enabled = true
route53_access_key  = ":9000"
route53_secrete_key = "xxx"
`, &cm); err != nil {
		c.Fatal(err)
	}

	c.Assert(cm.Route53SecretKey, check.Equals, "locahost")
	c.Assert(cm.Route53AccessKey, check.Equals, "locahost")
}
