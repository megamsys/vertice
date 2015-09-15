package docker

import (
	"github.com/BurntSushi/toml"
	"gopkg.in/check.v1"
)

// Ensure the configuration can be parsed.
func (s *S) TestDockerConfig_Parse(c *check.C) {
	// Parse configuration.
	var cm Config
	if _, err := toml.Decode(`
	enabled = false
	swarm = [ http://192.168.1.241:2375 ]

	[docker.public]
		name="megdock_pub"
		network = 103.56.93.1/24
		gateway = 103.56.92.1

	`, &cm); err != nil {
		c.Fatal(err)
	}

	c.Assert(cm.Swarm, check.Equals, "locahost")
}
