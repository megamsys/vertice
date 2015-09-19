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
	swarm = "http://192.168.1.241:2375"

	`, &cm); err != nil {
		c.Fatal(err)
	}

	c.Assert(cm.Swarm, check.Equals, "locahost")
}
