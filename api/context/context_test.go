package context

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/gorilla/context"
	"github.com/megamsys/megamd/auth"
	"gopkg.in/check.v1"
)

type S struct {
	token auth.Token
}

var _ = check.Suite(&S{})

func Test(t *testing.T) { check.TestingT(t) }

func (s *S) SetUpSuite(c *check.C) {
}

func (s *S) TestClear(c *check.C) {
	r, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, check.IsNil)
	context.Set(r, "my-key", "value")
	Clear(r)
	val := context.Get(r, "my-key")
	c.Assert(val, check.IsNil)
}

func (s *S) TestGetAuthToken(c *check.C) {
	r, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, check.IsNil)
	token := GetAuthToken(r)
	c.Assert(token, check.IsNil)
	SetAuthToken(r, s.token)
	token = GetAuthToken(r)
	c.Assert(token, check.Equals, s.token)
}

func (s *S) TestAddRequestError(c *check.C) {
	r, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, check.IsNil)
	err1 := errors.New("msg1")
	err2 := errors.New("msg2")
	myErr := GetRequestError(r)
	c.Assert(myErr, check.IsNil)
	AddRequestError(r, err1)
	myErr = GetRequestError(r)
	c.Assert(myErr, check.Equals, err1)
	AddRequestError(r, err2)
	otherErr := GetRequestError(r)
	c.Assert(otherErr.Error(), check.Equals, "msg2 Caused by: msg1")
}

func (s *S) TestSetDelayedHandler(c *check.C) {
	r, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, check.IsNil)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	val := GetDelayedHandler(r)
	c.Assert(val, check.IsNil)
	SetDelayedHandler(r, handler)
	val = GetDelayedHandler(r)
	v1 := reflect.ValueOf(val)
	v2 := reflect.ValueOf(handler)
	c.Assert(v1.Pointer(), check.Equals, v2.Pointer())
}
