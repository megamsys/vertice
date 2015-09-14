package docker

import (
	"github.com/BurntSushi/toml"
	"gopkg.in/check.v1"
)

// Ensure the configuration can be parsed.
func (s *S) TestConfig_Parse(c *check.C) {
	// Parse configuration.
	var c docker.Config
	if _, err := toml.Decode(`
	enabled = false
	master = [ http://192.168.1.241:2375 ]

	[docker.public]
		name="megdock_pub"
		network = 103.56.93.1/24
		gateway = 103.56.92.1

	`, &c); err != nil {
		t.Fatal(err)
	}

	c.Assert(c.Master, check.Equals, "locahost")
}
