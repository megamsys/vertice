package eventsd

import (
	"github.com/BurntSushi/toml"
	"gopkg.in/check.v1"
)

// Ensure the configuration can be parsed.
func (s *S) TestEventsConfig_Parse(c *check.C) {
	// Parse configuration.
	var cm Config
	if _, err := toml.Decode(`
	enabled = false
	`, &cm); err != nil {
		c.Fatal(err)
	}

	c.Assert(cm.Enabled, check.Equals, false)
}
