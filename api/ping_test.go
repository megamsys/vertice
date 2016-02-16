package api

import (
	/*	"net/http"
		"net/http/httptest"

		"github.com/BurntSushi/toml"
		"github.com/megamsys/vertice/meta"
	*/
	"gopkg.in/check.v1"
)

type HealthCheckSuite struct{}

var _ = check.Suite(&HealthCheckSuite{})

/*func (s *HealthCheckSuite) TestHealthCheck(c *check.C) {
	recorder := httptest.NewRecorder()
	var cm meta.Config
	if _, err := toml.Decode(`
	debug = true
	dir = "/var/lib/megam/vertice/meta"
	riak = ["localhost:8087"]
	api  = "http://localhost:9000"
	`, &cm); err != nil {
		c.Fatal(err)
	}
	cm.MkGlobal()
	request, err := http.NewRequest("GET", "/ping", nil)
	c.Assert(err, check.IsNil)
	healthcheck(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	//	c.Assert(recorder.Body.String(), check.Matches, `.*WORKING.*`)
}
*/
