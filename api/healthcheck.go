package api

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/megamsys/libgo/hc"
)

func healthcheck(w http.ResponseWriter, r *http.Request) error {
	fullHealthcheck(w, r)
	return nil
}

func fullHealthcheck(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	results := hc.Check()
	status := http.StatusOK
	for _, result := range results {
		fmt.Fprintf(&buf, "%s: %s (%s)\n", result.Name, result.Status, result.Duration)
		if result.Status != hc.HealthCheckOK {
			status = http.StatusInternalServerError
		}
	}
	w.WriteHeader(status)
	w.Write(buf.Bytes())
	//return c.Show(result)
}

/*
type ojak struct {
	LevelDB string
	Gw      string
	Nila    string
	Megamd  string
}

func (c *Ping) Show(result []byte) error {
	var o ojak
	err := json.Unmarshal(result, &o)
	if err != nil {
		return err
	}
	fmt.Fprintln(context.Stdout, &o)
	return nil
}
*/
