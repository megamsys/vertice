package httpd

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/pprof"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bmizerany/pat"
	"github.com/megamsys/megamd/meta"
)

type route struct {
	name        string
	method      string
	pattern     string
	handlerFunc interface{}
}

// Handler represents an HTTP handler for the Megamd server.
type Handler struct {
	Version string
	mux     *pat.PatternServeMux
}

// NewHandler returns a new instance of handler with routes.
func NewHandler() *Handler {
	h := &Handler{
		mux:            pat.New(),
		loggingEnabled: loggingEnabled,
	}

	h.SetRoutes([]route{
		route{ // Ping
			"ping",
			"GET", "/ping", true, h.servePing,
		},
	})

	return h
}

func (h *Handler) SetRoutes(routes []route) {
	for _, r := range routes {
		var handler http.Handler

		// This is a normal handler signature and does not require authorization
		if hf, ok := r.handlerFunc.(func(http.ResponseWriter, *http.Request)); ok {
			handler = http.HandlerFunc(hf)
		}

		handler = versionHeader(handler, h)
		h.mux.Add(r.method, r.pattern, handler)
	}
}

// ServeHTTP responds to HTTP request to the handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/debug/pprof") {
		switch r.URL.Path {
		case "/debug/pprof/cmdline":
			pprof.Cmdline(w, r)
		case "/debug/pprof/profile":
			pprof.Profile(w, r)
		case "/debug/pprof/symbol":
			pprof.Symbol(w, r)
		default:
			pprof.Index(w, r)
		}
		return
	}

	h.mux.ServeHTTP(w, r)
}

// servePing returns a simple response to let the client know the server is running.
func (h *Handler) servePing(w http.ResponseWriter, r *http.Request) {
	v := make(map[string]string)
	v["name"] = "megamd"
	v["version"] = "0.9"
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(&v)
}

// versionHeader takes a HTTP handler and returns a HTTP handler
// and adds the X-MEGAMD-VERSION header to outgoing responses.
func versionHeader(inner http.Handler, h *Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-MEGAMD-Version", h.Version)
		inner.ServeHTTP(w, r)
	})
}
