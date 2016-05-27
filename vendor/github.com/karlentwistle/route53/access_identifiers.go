package route53

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"time"
)

type AccessIdentifiers struct {
	AccessKey string
	SecretKey string
	time      time.Time
	endpoint  string
}

func (a *AccessIdentifiers) signature() (sha string) {
	if a.time.IsZero() {
		a.time = time.Now()
	}
	time := a.time.UTC().Format(time.ANSIC)
	hash := hmac.New(sha256.New, []byte(a.SecretKey))
	hash.Write([]byte(time))
	sha = base64.StdEncoding.EncodeToString(hash.Sum(nil))
	return
}

func (a *AccessIdentifiers) headers() http.Header {
	if a.time.IsZero() {
		a.time = time.Now()
	}
	signature := a.signature()
	h := http.Header{}
	h.Add("Date", a.time.UTC().Format(time.ANSIC))
	h.Add("Content-Type", "text/xml; charset=UTF-8")
	h.Add("X-Amzn-Authorization",
		"AWS3-HTTPS AWSAccessKeyId="+a.AccessKey+",Algorithm=HmacSHA256,Signature="+signature)
	return h
}

func (a *AccessIdentifiers) Zones() (h HostedZones) {
	var url string
	if a.endpoint == "" {
		url = awsURL + "?maxitems=100"
	} else {
		url = a.endpoint
	}
	res, err := a.zoneXML(url)
	if err == nil {
		return generateZones(res)
	}
	return h
}

func (a *AccessIdentifiers) zoneXML(url string) ([]byte, error) {
	resp, err := getBody(url, a.headers())
	if err == nil {
		return resp, err
	}
	return nil, err
}
