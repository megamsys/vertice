package metricsd

import (
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/check.v1"
)

// Ensure the configuration can be parsed.
func (s *S) TestDeploydConfig_Parse(c *check.C) {
	// Parse configuration.
	var cm Config
	if _, err := toml.Decode(`
		enabled = false
		collect_interval  = "10min"
`, &cm); err != nil {
		c.Fatal(err)
	}

	c.Assert(time.Duration(cm.CollectInterval), check.Equals, 15*time.Minute)
	c.Assert(cm.Enabled, check.Equals, true)

}
