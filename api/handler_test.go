package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/megamsys/libgo/errors"
	"github.com/megamsys/vertice/auth"
	"gopkg.in/check.v1"
)

type HandlerSuite struct {
	token auth.Token
}

var _ = check.Suite(&HandlerSuite{})

func (s *HandlerSuite) SetUpSuite(c *check.C) {
}

func errorHandler(w http.ResponseWriter, r *http.Request) error {
	return fmt.Errorf("some error")
}

func errorHandlerWriteHeader(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusBadGateway)
	return errorHandler(w, r)
}

func badRequestHandler(w http.ResponseWriter, r *http.Request) error {
	return &errors.HTTP{Code: http.StatusBadRequest, Message: "some error"}
}

func simpleHandler(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "success")
	return nil
}

func outputHandler(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text")
	output := "2012-06-05 17:03:36,887 WARNING ssl-hostname-verification is disabled for this environment"
	fmt.Fprint(w, output)
	return nil
}

func authorizedErrorHandler(w http.ResponseWriter, r *http.Request, t auth.Token) error {
	return errorHandler(w, r)
}

func authorizedErrorHandlerWriteHeader(w http.ResponseWriter, r *http.Request, t auth.Token) error {
	return errorHandlerWriteHeader(w, r)
}

func authorizedBadRequestHandler(w http.ResponseWriter, r *http.Request, t auth.Token) error {
	return badRequestHandler(w, r)
}

func authorizedSimpleHandler(w http.ResponseWriter, r *http.Request, t auth.Token) error {
	return simpleHandler(w, r)
}

func authorizedOutputHandler(w http.ResponseWriter, r *http.Request, t auth.Token) error {
	return outputHandler(w, r)
}

type recorder struct {
	*httptest.ResponseRecorder
	headerWrites int
}

func (r *recorder) WriteHeader(code int) {
	r.headerWrites++
	r.ResponseRecorder.WriteHeader(code)
}

/*
func (s *HandlerSuite) TestHandlerReturns500WhenInternalHandlerReturnsAnError(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/apps", nil)
	c.Assert(err, check.IsNil)
	RegisterHandler("/apps", "GET", Handler(errorHandler))
	defer resetHandlers()

}

func (s *HandlerSuite) TestHandlerDontCallWriteHeaderIfItHasAlreadyBeenCalled(c *check.C) {
	recorder := recorder{httptest.NewRecorder(), 0}
	request, err := http.NewRequest("GET", "/apps", nil)
	c.Assert(err, check.IsNil)
	RegisterHandler("/apps", "GET", Handler(errorHandlerWriteHeader))
	defer resetHandlers()

}

func (s *HandlerSuite) TestHandlerShouldPassAnHandlerWithoutError(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/apps", nil)
	c.Assert(err, check.IsNil)
	RegisterHandler("/apps", "GET", Handler(simpleHandler))
	defer resetHandlers()
}

func (s *HandlerSuite) TestAuthorizationRequiredHandlerShouldReturnUnauthorizedIfTheAuthorizationHeadIsNotPresent(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/apps", nil)
	c.Assert(err, check.IsNil)
	RegisterHandler("/apps", "GET", authorizationRequiredHandler(authorizedSimpleHandler))
	defer resetHandlers()
}

func (s *HandlerSuite) TestAuthorizationRequiredHandlerShouldReturnUnauthorizedIfTheTokenIsInvalid(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/apps", nil)
	c.Assert(err, check.IsNil)
	request.Header.Set("Authorization", "what the token?!")
	RegisterHandler("/apps", "GET", authorizationRequiredHandler(authorizedSimpleHandler))
	defer resetHandlers()
}

func (s *HandlerSuite) TestAuthorizationRequiredHandlerShouldReturnTheHandlerResultIfTheTokenIsOk(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/apps", nil)
	c.Assert(err, check.IsNil)
	request.Header.Set("Authorization", "bearer "+s.token.GetValue())
	RegisterHandler("/apps", "GET", authorizationRequiredHandler(authorizedSimpleHandler))
	defer resetHandlers()
}

func (s *HandlerSuite) TestAuthorizationRequiredHandlerShouldRespectTheHandlerStatusCode(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/apps", nil)
	c.Assert(err, check.IsNil)
	request.Header.Set("Authorization", "bearer "+s.token.GetValue())
	RegisterHandler("/apps", "GET", authorizationRequiredHandler(authorizedBadRequestHandler))
	defer resetHandlers()
}

func (s *HandlerSuite) TestAdminRequiredHandlerShouldReturnUnauthorizedIfTheAuthorizationHeadIsNotPresent(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/apps", nil)
	c.Assert(err, check.IsNil)
	RegisterHandler("/apps", "GET", authorizationRequiredHandler(authorizedSimpleHandler))
	defer resetHandlers()
}

func (s *HandlerSuite) TestAdminRequiredHandlerShouldReturnTheHandlerResultIfTheTokenIsOk(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/apps", nil)
	c.Assert(err, check.IsNil)
	request.Header.Set("Authorization", "bearer "+s.token.GetValue())
	RegisterHandler("/apps", "GET", AdminRequiredHandler(authorizedSimpleHandler))
	defer resetHandlers()
}
*/
