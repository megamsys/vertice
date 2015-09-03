package httpd

import (
	"github.com/BurntSushi/toml"
	"github.com/megamsys/megamd/services/httpd"
	"gopkg.in/check.v1"
)

// Ensure the configuration can be parsed.
func (s *S) TestConfig_Parse(c *check.C) {
	// Parse configuration.
	var c httpd.Config
	if _, err := toml.Decode(`
enabled = true
bind_address = ":8080"
`, &c); err != nil {
		t.Fatal(err)
	}

	c.Assert(c.BindAddress, check.Equals, "locahost")
}
