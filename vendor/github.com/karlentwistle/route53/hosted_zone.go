package route53

import (
	"encoding/xml"
	"net/http"
	"strings"
)

type CreateHostedZoneRequest struct {
	Name            string
	CallerReference string
	Comment         string `xml:"HostedZoneConfig>Comment"`
	Xmlns           string `xml:"xmlns,attr"`
}

func (hz *CreateHostedZoneRequest) XML() (s string, err error) {
	hz.Xmlns = `https://route53.amazonaws.com/doc/2012-12-12/`
	byteXML, err := xml.MarshalIndent(hz, "", `  `)
	if err != nil {
		return "", err
	}
	s = xml.Header + string(byteXML)
	return
}

func (hz *CreateHostedZoneRequest) Create(a AccessIdentifiers) (req *http.Response, err error) {
	postData, err := hz.XML()
	if err != nil {
		return nil, err
	}
	req, err = post(awsURL, postData, a.headers())
	return
}

type HostedZones struct {
	HostedZones []HostedZone `xml:"HostedZones>HostedZone"`
}

type HostedZone struct {
	Id              string `xml:"Id"`
	Name            string `xml:"Name"`
	CallerReference string `xml:"CallerReference"`
	RecordSetCount  int    `xml:"ResourceRecordSetCount"`
}

func (hz *HostedZone) HostedZoneId() string {
	s := strings.Split(hz.Id, "/")
	if len(s) == 3 {
		return s[2]
	}
	return ""
}

func (hz *HostedZone) ResourceRecordSets(a AccessIdentifiers) (r ResourceRecordSets, err error) {
	var url string
	if a.endpoint == "" {
		url = awsURL + "/" + hz.HostedZoneId() + "/rrset"
	} else {
		url = a.endpoint
	}

	resp, err := getBody(url, a.headers())
	if err == nil {
		return generateResourceRecordSet(resp), nil
	}
	return r, err
}
