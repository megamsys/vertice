package metrix

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	CPU_UNIT    = "cpu_unit"
	MEMORY_UNIT = "memory_unit"
	DISK_UNIT   = "disk_unit"
	CPU_COST    = "cpu_cost"
	MEMORY_COST = "memory_cost"
	CPU_COST_PER_HOUR = "cpu_cost_per_hour"
	RAM_COST_PER_HOUR = "ram_cost_per_hour"
	DISK_COST   = "disk_cost"

)

type MetricsMap map[string]int64

type Metrics []*Metric

func (list Metrics) Len() int {
	return len(list)
}

func (list Metrics) Swap(a, b int) {
	list[a], list[b] = list[b], list[a]
}

func flattenTags(tags map[string]string) string {
	out := []string{}
	for k, v := range tags {
		out = append(out, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(out)
	return strings.Join(out, " ")
}

func (p *Metrics) ToString() []string {
	swap := make([]string, 0)
	for _, j := range *p {
		b, _ := json.Marshal(j)
		swap = append(swap, string(b))
	}
	return swap
}

func parseStringToStruct(str string, data interface{}) error {
	if err := json.Unmarshal([]byte(str), data); err != nil {
		return err
	}
	return nil
}

func (m *Metrics) Totalcost(units map[string]string) string {

	//have to calculate metrics based on discount when flavour increases
	var cost float64
	defaultCpuUnit, _ := strconv.ParseFloat(units[CPU_UNIT], 64)
	defaultRamUnit, _ := strconv.ParseFloat(units[MEMORY_UNIT], 64)
	defaultDiskUnit, _ := strconv.ParseFloat(units[DISK_UNIT], 64)

	for _, in := range *m {
		consume, _ := strconv.ParseFloat(in.MetricValue, 64)
		unit, _ := strconv.ParseFloat(in.MetricUnits, 64)
		switch in.MetricName {
		case CPU_COST:
			cost = cost + (unit/defaultCpuUnit)*consume
		case MEMORY_COST:
			cost = cost + (unit/defaultRamUnit)*consume
		case DISK_COST:
			cost = cost + (unit/defaultDiskUnit)*consume
		}
	}
	return strconv.FormatFloat(cost/6, 'f', 6, 64)   //for 1 hr to 10min It should be customized
}
