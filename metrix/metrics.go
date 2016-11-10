package metrix

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)
var	defaultCpuUnit  float64 = 1
var	defaultDiskUnit  float64 = 24576
var	defaultRamUnit  float64  = 1024


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

func (m *Metrics) Totalcost() string {
	cost := 0.0
	for _, in := range *m {
		consume, _ := strconv.ParseFloat(in.MetricValue, 64)
		unit,_     := strconv.ParseFloat(in.MetricUnits, 64)
    switch in.MetricName {
    case CPU_COST:
        cost = cost + (unit / defaultCpuUnit ) * consume
		case MEMORY_COST:
        cost = cost + ( unit / defaultRamUnit ) * consume
		case DISK_COST:
        cost = cost + (unit / defaultDiskUnit) * consume
    }
	}
	return strconv.FormatFloat(cost, 'f', 3, 64)
}
