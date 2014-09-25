package app

import (
	"errors"
	"gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type S struct{}

var _ = check.Suite(&S{})

func (s *S) TestAppLifecycleError(c *check.C) {
	e := AppLifecycleError{app: "myapp", Err: errors.New("failure in app")}
	expected := `gulpd failed to apply the lifecle to the app "myapp": failure in app`
	c.Assert(e.Error(), check.Equals, expected)
}
