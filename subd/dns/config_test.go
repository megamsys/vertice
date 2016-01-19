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
access_key  = ":9000"
secret_key = "xxx"
`, &cm); err != nil {
		c.Fatal(err)
	}

	c.Assert(cm.SecretKey, check.Equals, "xxx")
	c.Assert(cm.AccessKey, check.Equals, ":9000")
}
