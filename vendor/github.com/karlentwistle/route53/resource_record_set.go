package route53

import (
	"encoding/xml"
	"net/http"
)

type ChangeResourceRecordSetsRequest struct {
	ZoneID  string   `xml:"-"`
	Comment string   `xml:"ChangeBatch>Comment"`
	Changes []Change `xml:"ChangeBatch>Changes>Change"`
	Xmlns   string   `xml:"xmlns,attr"`
}

type Change struct {
	Action string
	Name   string `xml:"ResourceRecordSet>Name"`
	Type   string `xml:"ResourceRecordSet>Type"`
	TTL    int    `xml:"ResourceRecordSet>TTL"`
	Value  string `xml:"ResourceRecordSet>ResourceRecords>ResourceRecord>Value"`
}

func (c *ChangeResourceRecordSetsRequest) XML() (s string, err error) {
	c.Xmlns = `https://route53.amazonaws.com/doc/2012-12-12/`
	byteXML, err := xml.MarshalIndent(c, "", `   `)
	if err != nil {
		return "", err
	}
	s = xml.Header + string(byteXML)
	return
}

func (c *ChangeResourceRecordSetsRequest) Create(a AccessIdentifiers) (req *http.Response, err error) {
	postData, err := c.XML()
	if err != nil {
		return nil, err
	}
	url := awsURL + `/` + c.ZoneID + `/rrset`
	req, err = post(url, postData, a.headers())
	return
}

type ResourceRecordSets struct {
	ResourceRecordSets []ResourceRecordSet `xml:"ResourceRecordSets>ResourceRecordSet"`
}

type ResourceRecordSet struct {
	Name  string   `xml:"Name"`
	Type  string   `xml:"Type"`
	TTL   int      `xml:"TTL"`
	Value []string `xml:"ResourceRecords>ResourceRecord>Value"`
}
