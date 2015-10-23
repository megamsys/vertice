package run

/*
// Ensure the configuration can be parsed.
func (s *S) TestConfig_Parse(c *check.C) {
	var cm *Config = NewConfig()
	path := cm.Meta.Dir + "/megamd.conf"
	c.Assert((len(strings.TrimSpace(path)) > 0), check.Equals, true)
	if _, err := toml.DecodeFile(path, cm); err != nil {
		fmt.Println(err.Error())
	}

	c.Assert(cm, check.NotNil)
	c.Assert(cm.Meta.Hostname, check.Equals, "localhost")
	c.Assert(cm.Meta.Riak, check.DeepEquals, []string{"localhost:8087"})
	c.Assert(cm.Meta.Api, check.Equals, "https://api.megam.io/v2")
	c.Assert(cm.Deployd.OneEndPoint, check.Equals, "http://192.168.1.100:3030/RPC2")
	c.Assert(cm.Deployd.OneUserid, check.Equals, "oneadmin")
}
*/
