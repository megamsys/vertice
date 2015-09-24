package run

import (
	"github.com/BurntSushi/toml"
	"gopkg.in/check.v1"
)


// Ensure the configuration can be parsed.
func (s *S) TestConfig_Parse(c *check.C) {
	// Parse configuration.
	var cm Config
	if _, err := toml.Decode(`
		[meta]
			debug = true
			hostname = "localhost"
			bind_address = ":9999"
			dir = "/var/lib/megam/megamd/meta"
			riak = "192.168.1.100:8087"
			api  = "https://api.megam.io/v2"
			amqp = "amqp://guest:guest@192.168.1.100:5672/"

		[deployd]
			one_endpoint = "http://192.168.1.100:3030/xmlrpc2"
			one_userid = "oneadmin"
			one_password =  "password"

		[http]
	    enabled = true
	    bind-address = "localhost:6666"

	  [docker]
	    enabled = false
	    master = [ http://192.168.1.241:2375 ]

    [docker.public]
       name="megdock_pub"
       network = 103.56.93.1/24
       gateway = 103.56.92.1

   [dns]
	    enabled = true
	    route53_access_key = "accesskey"
	    route53_secret_key = "secretkey"
`, &cm); err != nil {
		c.Fatal(err)
	}

	c.Assert(cm.Meta.Hostname, check.Equals, "localhost")
	c.Assert(cm.Meta.Riak, check.Equals, "192.168.1.100:8087")
	c.Assert(cm.Meta.Api, check.Equals, "https://api.megam.io/v2")
	c.Assert(cm.Deployd.OneEndPoint, check.Equals, "http://192.168.1.100:3030/xmlrpc2")
	c.Assert(cm.Deployd.OneUserid, check.Equals, "onedmin")

}
