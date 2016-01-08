package run

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/check.v1"
)

// Ensure the configuration can be parsed.
func (s *S) TestConfig_Parse(c *check.C) {
	var cm *Config = NewConfig()
	path := cm.Meta.Dir + "/megamd.conf"

	c.Assert((len(strings.TrimSpace(path)) > 0), check.Equals, true)
	if _, err := toml.DecodeFile(path, cm); err != nil {
		fmt.Println(err.Error())
	}

	c.Assert(cm, check.NotNil)
	c.Assert(cm.Meta.Riak, check.DeepEquals, []string{"103.56.92.7:8087"})
	c.Assert(cm.Meta.NSQd, check.DeepEquals, []string{"103.56.92.7:4150"})
	c.Assert(cm.Deployd.OneUserid, check.Equals, "oneadmin")

}
