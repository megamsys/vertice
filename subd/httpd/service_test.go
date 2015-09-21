package httpd

import (
	"fmt"
	"net/http"
	//"net/http/httptest"

	"github.com/megamsys/megamd/auth"
//	"gopkg.in/check.v1"
)

func authorizedMegamHandler(w http.ResponseWriter, r *http.Request, t auth.Token) error {
	fmt.Fprint(w, r.Method)
	return nil
}

/*
func (s *S) TestRegisterHandlerMakesHandlerAvailableViaGet(c *check.C) {
	RegisterHandler("/foo/bar", "GET", authorizationRequiredHandler(authorizedTsuruHandler))
	defer resetHandlers()
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://example.com/foo/bar", nil)
	req.Header.Set("Authorization", "bearer "+s.token.GetValue())
	c.Assert(err, check.IsNil)
	m := RunServer(true)
	m.ServeHTTP(rec, req)
	b, err := ioutil.ReadAll(rec.Body)
	c.Assert(err, check.IsNil)
	c.Assert("GET", check.Equals, string(b))

}

func (s *S) TestRegisterHandlerMakesHandlerAvailableViaPost(c *check.C) {
	RegisterHandler("/foo/bar", "POST", authorizationRequiredHandler(authorizedTsuruHandler))
	defer resetHandlers()
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "http://example.com/foo/bar", nil)
	req.Header.Set("Authorization", "bearer "+s.token.GetValue())
	c.Assert(err, check.IsNil)
	m := RunServer(true)
	m.ServeHTTP(rec, req)
	b, err := ioutil.ReadAll(rec.Body)
	c.Assert(err, check.IsNil)
	c.Assert("POST", check.Equals, string(b))
}

func (s *S) TestRegisterHandlerMakesHandlerAvailableViaPut(c *check.C) {
	RegisterHandler("/foo/bar", "PUT", authorizationRequiredHandler(authorizedTsuruHandler))
	defer resetHandlers()
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "http://example.com/foo/bar", nil)
	req.Header.Set("Authorization", "bearer "+s.token.GetValue())
	c.Assert(err, check.IsNil)
	m := RunServer(true)
	m.ServeHTTP(rec, req)
	b, err := ioutil.ReadAll(rec.Body)
	c.Assert(err, check.IsNil)
	c.Assert("PUT", check.Equals, string(b))
}

func (s *S) TestRegisterHandlerMakesHandlerAvailableViaDelete(c *check.C) {
	RegisterHandler("/foo/bar", "DELETE", authorizationRequiredHandler(authorizedTsuruHandler))
	defer resetHandlers()
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "http://example.com/foo/bar", nil)
	req.Header.Set("Authorization", "bearer "+s.token.GetValue())
	c.Assert(err, check.IsNil)
	m := RunServer(true)
	m.ServeHTTP(rec, req)
	b, err := ioutil.ReadAll(rec.Body)
	c.Assert(err, check.IsNil)
	c.Assert("DELETE", check.Equals, string(b))
}
*/
