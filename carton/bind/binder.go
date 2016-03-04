/*
** Copyright [2013-2016] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */
package bind

import (
	"encoding/json"
	"fmt"
	"github.com/megamsys/libgo/os"
	"runtime"
	"strings"
)

// EnvVar represents a environment variable for a carton.
type EnvVar struct {
	Name     string
	Value    string
	Endpoint string
}

func (e *EnvVar) String() string {
	return fmt.Sprintf("%s=%s", e.Name, e.Value)
}

type EnvVars []EnvVar

func (en EnvVars) WrapForInitds() string {
	var envs = ""
	for _, de := range en {
		envs += wrapForInitdservice(de.Name, de.Value)
	}
	return envs
}

func wrapForInitdservice(key string, value string) string {
	osh := os.HostOS()
	switch runtime.GOOS {
	case "linux":
		if osh != os.Ubuntu {
			return key + "=" + value //systemd
		}
	default:
		return "initctl set-env " + key + "=" + value + "\n"
	}
	return "initctl set-env " + key + "=" + value + "\n"
}

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
