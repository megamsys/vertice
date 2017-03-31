package api

import (
	"fmt"
	"github.com/megamsys/vertice/api/context"
	"net/http"
)

type Handler func(http.ResponseWriter, *http.Request) error

func (fn Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context.AddRequestError(r, fn(w, r))
}
