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
	"github.com/megamsys/vertice/api"
	"time"
)

const SENSORSBUCKET = "sensors"

type Sensor struct {
	Id                   string  `json:"id" cql:"id"`
	AccountId            string  `json:"account_id" cql:"account_id"`
	SensorType           string  `json:"sensor_type" cql:"sensor_type"`
	AssemblyId           string  `json:"assembly_id" cql:"assembly_id"`
	AssemblyName         string  `json:"assembly_name" cql:"assembly_name"`
	AssembliesId         string  `json:"assemblies_id" cql:"assemblies_id"`
	Node                 string  `json:"node" cql:"node"`
	System               string  `json:"system" cql:"system"`
	Status               string  `json:"status" cql:"status"`
	Source               string  `json:"source" cql:"source"`
	Message              string  `json:"message" cql:"message"`
	AuditPeriodBeginning string  `json:"audit_period_beginning" cql:"audit_period_beginning"`
	AuditPeriodEnding    string  `json:"audit_period_ending" cql:"audit_period_ending"`
	AuditPeriodDelta     string  `json:"audit_period_delta" cql:"audit_period_delta"`
	Metrics              Metrics `json:"metrics" cql:"metrics"`
	CreatedAt            time.Time  `json:"created_at" cql:"created_at"`
}


func (s *Sensor) String() string {
	insp, _ := json.Marshal(s)
	o, _ := api.PP(insp)
	return string(o)
}

func NewSensor(senstype string) *Sensor {
	s := &Sensor{
		Id:         api.Uid(""),
		SensorType: senstype,
	}
	return s
}

func (s *Sensor) addMetric(key string, value string, unit string, mtype string) {
	s.newMetric(
		&Metric{
			MetricName:  key,
			MetricValue: value,
			MetricUnits: unit,
			MetricType:  mtype,
		})
}

func (s *Sensor) newMetric(me *Metric) {
	s.Metrics = append(s.Metrics, me)
}

func (s *Sensor) isBillable() bool {
	return s.AccountId != "" && s.AssemblyName != ""
}

func (s *Sensor) WriteIt() {
}
