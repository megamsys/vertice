package meta

import (
	"github.com/BurntSushi/toml"
	"gopkg.in/check.v1"
)

// Ensure the configuration can be parsed.
func (s *S) TestMetaConfig_Parse(c *check.C) {
	// Parse configuration.
	var cm Config
	if _, err := toml.Decode(`
debug = true
dir = "/var/lib/megam/megamd/meta"
riak = ["192.168.1.100:8087"]
api  = "http://localhost:9000"
amqp = "amqp://guest:guest@192.168.1.100:5672/"
`, &cm); err != nil {
		c.Fatal(err)
	}

	c.Assert(cm.Api, check.Equals, "http://localhost:9000")

}
