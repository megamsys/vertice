package docker

import (
	"github.com/BurntSushi/toml"
	"gopkg.in/check.v1"
)

// Ensure the configuration can be parsed.
func (s *S) TestDockerBrigeConfig_Parse(c *check.C) {
	// Parse configuration.
	var cm Bridges
	if _, err := toml.Decode(`
	[bridges]

    [bridges.public]
		  name = "megdock_pub"
		  network = "103.56.93.1/24"
		  gateway = "103.56.92.1"

    [bridges.private]
      name = "megdock_private"
		  network = "192.168.1.128/24"
		  gateway = "192.168.1.1"

	`, &cm); err != nil {
		c.Fatal(err)
	}
	c.Assert(cm, check.HasLen, 2)
}
