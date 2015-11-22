package api

import (
	"encoding/json"
	"net/http"

	"github.com/megamsys/libgo/hc"
	_ "github.com/megamsys/megamd/hc"
)

func healthcheck(w http.ResponseWriter, r *http.Request) error {
	fullHealthcheck(w, r)
	return nil
}

func fullHealthcheck(w http.ResponseWriter, r *http.Request) {
	results := hc.Check()
	status := http.StatusOK
	for _, result := range results {
		if result.Status != hc.HealthCheckOK {
			status = http.StatusInternalServerError
		}
	}
	data, _ := json.MarshalIndent(results, "", "  ")
	w.WriteHeader(status)
	w.Write(data)
}
