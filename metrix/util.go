package metrix

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

func AbortWith(message string) {
	log.Debugf("ERROR: %s", message)
	flag.PrintDefaults()
	os.Exit(1)
}

func FetchURL(url string) (b []byte, e error) {
	var rsp *http.Response
	rsp, e = http.Get(url)
	if e != nil {
		return
	}
	defer rsp.Body.Close()
	b, e = ioutil.ReadAll(rsp.Body)
	return
}
