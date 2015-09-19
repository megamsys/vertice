package run

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"gopkg.in/check.v1"
	"os/user"
)

// Ensure the configuration can be parsed.
func (s *S) TestConfig_Parse(c *check.C) {
  var cm Config
	u, _ := user.Current()
	if _, err := toml.DecodeFile(u.HomeDir+"/code/megam/go/src/github.com/megamsys/megamd/conf/megamd.conf", &cm); err != nil {
		fmt.Println(err.Error())
	}

	c.Assert(cm, check.NotNil)
	c.Assert(cm.Meta.Hostname, check.Equals, "localhost")
	c.Assert(cm.Meta.Riak, check.DeepEquals, []string{"localhost:8087"})
	c.Assert(cm.Meta.Api, check.Equals, "https://api.megam.io/v2")
	c.Assert(cm.Deployd.OneEndPoint, check.Equals, "http://192.168.1.100:3030/xmlrpc2")
	c.Assert(cm.Deployd.OneUserid, check.Equals, "oneadmin")
}
