package hc

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/megamsys/libgo/hc"
	"github.com/megamsys/megamd/meta"
)

var httpRegexp = regexp.MustCompile(`^http?://`)

func init() {
	hc.AddChecker("gateway", healthCheckGW)
}

func healthCheckGW() error {
	api := meta.MC.Api
	if api == "" {
		return hc.ErrDisabledComponent
	}

	if !httpRegexp.MatchString(api) {
		api = "http://" + api
	}
	v2URL := strings.TrimRight(api, "/")
	_, err := http.Get(v2URL)
	if err != nil {
		return err
	}
	return nil
}
