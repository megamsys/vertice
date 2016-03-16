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

	"github.com/megamsys/vertice/api"
)

const SENSORSBUCKET = "sensors"

type Sensor struct {
	Id         string   `json:"id" 			cql:"id"`
	AccountsId string   `json:"account_id" 	cql:"accounts_id"`
	SensorType string   `json:"sensor_type" cql:"sensor_type"`
	Payload    *Payload `json:"payload" 	cql:"payload"`
	CreatedAt  string   `json:"created_at" 	cql:"created_at"`
}

type SensorScylla struct {
	Id         string `json:"id"			cql:"id"`
	AccountsId string `json:"accounts_id"	cql:"accounts_id"`
	SensorType string `json:"sensor_type"	cql:"sensor_type"`
	Payload    string `json:"payload"		cql:"payload"`
	CreatedAt  string `json:"created_at"	cql:"created_at"`
}

func (s *Sensor) String() string {
	insp, _ := json.Marshal(s)
	o, _ := api.PP(insp)
	return string(o)
}

type Payload struct {
	AssemblyId   string `json:"assembly_id"				cql:"assembly_id"`
	AssemblyName string `json:"assembly_name"			cql:"assembly_name"`
	AssembliesId string `json:"assemblies_id"			cql:"assemblies_id"`
	Node         string `json:"node"					cql:"node"`
	System       string `json:"system"					cql:"system"`
	Status       string `json:"status"					cql:"status"`
	Source       string `json:"source"					cql:"source"`
	Message      string `json:"message"					cql:"message"`
	BeginAudit   string `json:"audit_period_beginning"	cql:"audit_period_beginning"`
	EndAudit     string `json:"audit_period_ending"		cql:"audit_period_ending"`
	DeltaAudit   string `json:"audit_period_delta"		cql:"audit_period_delta"`

	Metrics []*Metric `json:"metrics"					cql:"metrics"`
}

func NewSensor(senstype string) *Sensor {
	s := &Sensor{
		Id:         api.Uid(""),
		SensorType: senstype,
		CreatedAt:  time.Now().Local().Format(time.RFC822),
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

func (s *Sensor) ParseScyllaformat() (SensorScylla, error) {
	b, err := json.Marshal(s.Payload)
	return SensorScylla{
		Id:         s.Id,
		AccountsId: s.AccountsId,
		SensorType: s.SensorType,
		Payload:    string(b),
		CreatedAt:  s.CreatedAt,
	}, err
}

func (pa *Payload) newMetric(me *Metric) {
	pa.Metrics = append(pa.Metrics, me)
}

func (s *Sensor) WriteIt() {
}
