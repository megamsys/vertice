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
		[events]
	    enabled = false

	  [mailgun]
	    api_key = "temp"
	    domain  = "ojamail.megambox.com"

	  [infobip]
	    username = "info_username"
	    password = "info_pw"
	    api_key  = "info_apiky"
	    application_id = "info_apiid"
	    message_id = "info_msgid"

	  [slack]
	    api_key = "temp"
	    room  = "bugs"

	  [bill]
	    api_key = "whmcs"
	`, &cm); err != nil {
		c.Fatal(err)
	}

	c.Assert(cm.Enabled, check.Equals, false)
	c.Assert(cm.Mailgun.ApiKey, check.Equals, "temp")
	c.Assert(cm.Infobip.Username, check.Equals, "info_username")
	c.Assert(cm.Slack.ApiKey, check.Equals, "temp")
}
