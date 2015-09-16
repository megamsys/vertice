package httpd

import (
	"gopkg.in/check.v1"

)

// NewHandler represents a test wrapper for httpd.Handler.
type THh struct {
	th *Handler
}

// NewHandler returns a new instance of Handler.
func NewTHh() *THh {
	t := &THh{
		th: NewHandler(),
	}
	t.th.Version = "0.0.0"
	return t
}

func (s *S) SetUpSuite(c *check.C) {
	h := NewTHh()
	c.Assert(h, check.NotNil)
}
