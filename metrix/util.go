package metrix

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	log "github.com/Sirupsen/logrus"
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
