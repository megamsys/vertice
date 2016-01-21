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
package metrix

import (
	"encoding/json"
	"time"

	"github.com/megamsys/megamd/api"
)

const SENSORSBUCKET = "sensors"

type Sensor struct {
	Id         string   `json:"id"`
	AccountsId string   `json:"accounts" riak:"index"`
	Type       string   `json:"type"`
	Payload    *Payload `json:"payload"`
	CreatedAt  string   `json:"created_at"`
}

func (s *Sensor) String() string {
	insp, _ := json.Marshal(s)
	o, _ := api.PP(insp)
	return string(o)
}

type Payload struct {
	AssemblyId   string `json:"assembly_id"`
	AssemblyName string `json:"assembly_name"`
	AssembliesId string `json:"assemblies_id"`
	Node         string `json:"node"`
	System       string `json:"system"`
	Status       string `json:"status"`
	Source       string `json:"source"`
	Message      string `json:"message"`
	BeginAudit   string `json:"audit_period_beginning"`
	EndAudit     string `json:"audit_period_ending"`
	DeltaAudit   string `json:"audit_period_delta"`

	Metrics []*Metric `json:"metrics"`
}

func NewSensor(senstype string) *Sensor {
	s := &Sensor{
		Id:        api.Uid(""),
		Type:      senstype,
		CreatedAt: time.Now().Local().Format(time.RFC822),
	}
	return s
}

func (s *Sensor) addPayload(pa *Payload) {
	s.Payload = pa
}

func (s *Sensor) addMetric(key string, value string, unit string, mtype string) {
	s.Payload.newMetric(
		&Metric{
			Key:   key,
			Value: value,
			Units: unit,
			Type:  mtype,
		})
}

func (pa *Payload) newMetric(me *Metric) {
	pa.Metrics = append(pa.Metrics, me)
}

func (s *Sensor) WriteIt() {
}
