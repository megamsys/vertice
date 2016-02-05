package hc

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/megamsys/libgo/hc"
	"github.com/megamsys/vertice/meta"
)

const PORT = "8098"
const LEVEL_DB = "riak_kv_eleveldb_backend"
const LEVEL_DB_NOT_SET = "Missing storage_backend: leveldb in /etc/riak/riak.conf"

func init() {
	hc.AddChecker("leveldb", healthCheckLdb)
}

func healthCheckLdb() (interface{}, error) {
	apia := meta.MC.Riak
	if len(apia) <= 0 {
		return nil, hc.ErrDisabledComponent
	}
	api := apia[0]
	api = strings.Split(api, ":")[0] + ":" + PORT
	if !httpRegexp.MatchString(api) {
		api = "http://" + api
	}

	api = strings.TrimRight(api, "/")
	vRURL := api + "/stats"
	response, err := http.Get(vRURL)

	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(response.Body)

	if strings.Contains(string(body), LEVEL_DB) {
		return LEVEL_DB, nil
	}
	return LEVEL_DB_NOT_SET, nil
}
