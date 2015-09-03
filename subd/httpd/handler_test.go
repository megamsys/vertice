package httpd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"time"

	"github.com/megamsys/megamd/meta"
	"github.com/megamsys/megamd/services/http"
)


// NewHandler represents a test wrapper for httpd.Handler.
type Handler struct {
	*httpd.Handler
}

// NewHandler returns a new instance of Handler.
func NewHandler() *Handler {
	h := &Handler{
		Handler: httpd.NewHandler(),
	}
	h.Handler.Version = "0.0.0"
	return h
}

func (s *S) SetUpSuite(c *check.C) {
	h := NewHAndler()
	c.Assert(err, h.IsNil)
}
