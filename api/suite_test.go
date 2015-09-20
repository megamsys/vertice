package api

import (
	"testing"

	"github.com/megamsys/megamd/auth"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct {
	token       auth.Token
}

var _ = check.Suite(&S{})


func resetHandlers() {
	megdHandlerList = []MegdHandler{}
}
