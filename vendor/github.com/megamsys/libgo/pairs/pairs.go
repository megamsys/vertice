package pairs

import (
	"encoding/json"
	"strings"
)

type JsonPair struct {
	K string `json:"key" cql:"key"`
	V string `json:"value" cql:"value"`
}

type JsonPairs []*JsonPair

func NewJsonPair(k string, v string) *JsonPair {
	return &JsonPair{
		K: k,
		V: v,
	}
}

//match for a value in the JSONPair and send the value
func (p *JsonPairs) Match(k string) string {
	for _, j := range *p {
		if j.K == k {
			return j.V
		}
	}
	return ""
}

func (p *JsonPairs) ToMap() map[string]string {
	jm := make(map[string]string)
	for _, j := range *p {
		jm[j.K] = j.V
	}
	return jm
}

func (p *JsonPairs) ToString() []string {
	swap := make([]string, 0)
	for _, j := range *p {
		b, _ := json.Marshal(j)
		swap = append(swap, string(b))
	}
	return swap
}

//Delete old keys and update them with the new values
func (p *JsonPairs) NukeAndSet(m map[string][]string) {
	swap := make(JsonPairs, 0)
	for _, j := range *p { //j is value
		exists := false
		for k, _ := range m { //k is key
			if strings.Compare(j.K, k) == 0 {
				exists = true
				break
			}
		}
		if !exists {
			swap = append(swap, j)
		}
	}
	for mkey, mvals := range m {
		for _, mval := range mvals {
			swap = append(swap, NewJsonPair(mkey, mval))
		}
	}
	*p = swap
}
