package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/megamsys/libgo/errors"
	"github.com/megamsys/libgo/io"
	"github.com/megamsys/vertice/api/context"
	"gopkg.in/check.v1"
)

type handlerLog struct {
	w        http.ResponseWriter
	r        *http.Request
	called   bool
	sleep    time.Duration
	response int
}

func doHandler() (http.HandlerFunc, *handlerLog) {
	h := &handlerLog{}
	return func(w http.ResponseWriter, r *http.Request) {
		if h.sleep != 0 {
			time.Sleep(h.sleep)
		}
		h.called = true
		h.w = w
		h.r = r
		if h.response != 0 {
			w.WriteHeader(h.response)
		}
	}, h
}

func (s *S) TestContextClearerMiddleware(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, check.IsNil)
	context.AddRequestError(request, fmt.Errorf("Some Error"))
	h, log := doHandler()
	contextClearerMiddleware(recorder, request, h)
	c.Assert(log.called, check.Equals, true)
	contErr := context.GetRequestError(request)
	c.Assert(contErr, check.IsNil)
}

func (s *S) TestFlushingWriterMiddleware(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, check.IsNil)
	h, log := doHandler()
	flushingWriterMiddleware(recorder, request, h)
	c.Assert(log.called, check.Equals, true)
	_, ok := log.w.(*io.FlushingWriter)
	c.Assert(ok, check.Equals, true)
}

func (s *S) TestErrorHandlingMiddlewareWithoutError(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, check.IsNil)
	h, log := doHandler()
	errorHandlingMiddleware(recorder, request, h)
	c.Assert(log.called, check.Equals, true)
	c.Assert(recorder.Code, check.Equals, 200)
}

func (s *S) TestErrorHandlingMiddlewareWithError(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, check.IsNil)
	h, log := doHandler()
	context.AddRequestError(request, fmt.Errorf("something"))
	errorHandlingMiddleware(recorder, request, h)
	c.Assert(log.called, check.Equals, true)
	c.Assert(recorder.Code, check.Equals, 500)
}

func (s *S) TestErrorHandlingMiddlewareWithHTTPError(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, check.IsNil)
	h, log := doHandler()
	context.AddRequestError(request, &errors.HTTP{Code: 403, Message: "other msg"})
	errorHandlingMiddleware(recorder, request, h)
	c.Assert(log.called, check.Equals, true)
	c.Assert(recorder.Code, check.Equals, 403)
}

func (s *S) TestAuthTokenMiddlewareWithoutToken(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, check.IsNil)
	h, log := doHandler()
	authTokenMiddleware(recorder, request, h)
	c.Assert(log.called, check.Equals, true)
	t := context.GetAuthToken(request)
	c.Assert(t, check.IsNil)
}

func (s *S) TestAuthTokenMiddlewareWithToken(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, check.IsNil)
	request.Header.Set("Authorization", "bearer "+s.token.GetValue())
	h, log := doHandler()
	authTokenMiddleware(recorder, request, h)
	c.Assert(log.called, check.Equals, true)
	t := context.GetAuthToken(request)
	c.Assert(t, check.NotNil)
	c.Assert(t.GetValue(), check.Equals, s.token.GetValue())
	c.Assert(t.GetUserName(), check.Equals, s.token.GetUserName())
}

func (s *S) TestRunDelayedHandlerWithoutHandler(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, check.IsNil)
	runDelayedHandler(recorder, request)
}

func (s *S) TestRunDelayedHandlerWithHandler(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, check.IsNil)
	h, log := doHandler()
	context.SetDelayedHandler(request, h)
	runDelayedHandler(recorder, request)
	c.Assert(log.called, check.Equals, true)
}
