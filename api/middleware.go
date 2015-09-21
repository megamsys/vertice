package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/megamsys/libgo/errors"
	"github.com/megamsys/libgo/io"
	"github.com/megamsys/megamd/api/context"
	"github.com/megamsys/megamd/auth"
)

func validate(token string, r *http.Request) (auth.Token, error) {
	t, err := Auth(token)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func contextClearerMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer context.Clear(r)
	next(w, r)
}

func flushingWriterMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if r.Body != nil {
			r.Body.Close()
		}
	}()
	fw := io.FlushingWriter{ResponseWriter: w}
	next(&fw, r)
}

func errorHandlingMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(w, r)
	err := context.GetRequestError(r)
	if err != nil {
		code := http.StatusInternalServerError
		if e, ok := err.(*errors.HTTP); ok {
			code = e.Code
		}
		flushing, ok := w.(*io.FlushingWriter)
		if ok && flushing.Wrote() {
			fmt.Fprintln(w, err)
		} else {
			http.Error(w, err.Error(), code)
		}
		log.Errorf("failure running HTTP request %s %s (%d): %s", r.Method, r.URL.Path, code, err)
	}
}

func authTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token := r.Header.Get("Authorization")
	if token != "" {
		t, err := validate(token, r)
		if err != nil {
			if err != auth.ErrInvalidToken {
				context.AddRequestError(r, err)
				return
			}
			log.Debugf("Ignored invalid token for %s: %s", r.URL.Path, err.Error())
		} else {
			context.SetAuthToken(r, t)
		}
	}
	next(w, r)
}

func runDelayedHandler(w http.ResponseWriter, r *http.Request) {
	h := context.GetDelayedHandler(r)
	if h != nil {
		h.ServeHTTP(w, r)
	}
}

type loggerMiddleware struct {
	logger *log.Logger
}

func newLoggerMiddleware() *loggerMiddleware {
	l := log.New()
	l.Out = os.Stdout

	return &loggerMiddleware{
		logger: l,
	}
}

func (l *loggerMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(rw, r)
	duration := time.Since(start)
	res := rw.(negroni.ResponseWriter)
	nowFormatted := time.Now().Format(time.RFC3339Nano)
	l.logger.Debugf("%s %s %s %d in %0.6fms", nowFormatted, r.Method, r.URL.Path, res.Status(), float64(duration)/float64(time.Millisecond))
}
