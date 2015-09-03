package meta

import (
	"time"

	"github.com/BurntSushi/toml"
	"github.com/megamsys/megamd/meta"
	"gopkg.in/check.v1"
)

// Ensure the configuration can be parsed.
func (s *S) TestConfig_Parse(c *check.C) {
	// Parse configuration.
	var c run.Config
	if _, err := toml.Decode(`
		[meta]
			debug = true
			hostname = "localhost"
			bind_address = ":9999"
			dir = "/var/lib/megam/megamd/meta"
			riak = "192.168.1.100:8087"
			api  = "https://api.megam.io/v2"
			amqp = "amqp://guest:guest@192.168.1.100:5672/"

`, &c); err != nil {
		t.Fatal(err)
	}

	c.Assert(c.Meta.hostname, check.Equals, "locahost")
	c.Assert(c.Meta.riak, check.Equals, "locahost")
	c.Assert(c.Meta.api, check.Equals, "locahost")

}
