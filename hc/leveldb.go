package hc

import (
	"net/http"
	"strings"

	"github.com/megamsys/libgo/hc"
	"github.com/megamsys/megamd/meta"
)

func init() {
	hc.AddChecker("leveldb", healthCheckLdb)
}

func healthCheckLdb() error {
	apia := meta.MC.Riak
	if len(apia) <= 0 {
		return hc.ErrDisabledComponent
	}
	api := apia[0]
	if !httpRegexp.MatchString(api) {
		api = "http://" + api
	}
	api = strings.TrimRight(api, "/")
	vRURL := api + "stats"
	_, err := http.Get(vRURL)
	if err != nil {
		return err
	}
	return nil
}
