package hc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/megamsys/libgo/hc"
	"github.com/megamsys/vertice/meta"
)

const GATEWAY_ERROR = "Missing initial info: to fix, run curl %s"

var httpRegexp = regexp.MustCompile(`^http?://`)

func init() {
	hc.AddChecker("gateway", healthCheckGW)
}

func healthCheckGW() (interface{}, error) {
	api := meta.MC.Api
	if api == "" {
		return nil, hc.ErrDisabledComponent
	}

	if !httpRegexp.MatchString(api) {
		api = "http://" + api
	}

	v2URL := strings.TrimRight(api, "/")
	response, err := http.Get(v2URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	return mkGateway(body, api), nil

}

func mkGateway(data []byte, api string) Gateway {
	g := Gateway{}
	if err := json.Unmarshal(data, &g); err != nil {
		g.Status = map[string]string{"error": fmt.Sprintf(GATEWAY_ERROR, api)}
		return g
	}
	return g
}

type Gateway struct {
	Status  map[string]string `json:"status"`
	Runtime map[string]string `json:"runtime"`
	Loaded  map[string]string `json:"loaded"`
}
