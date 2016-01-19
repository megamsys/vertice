package metrix

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

type Metric struct {
	Key   string
	Value string
	Tags  map[string]string `json:"Tags,omitempty"`
}

func (m *Metric) String() string {
	if d, err := yaml.Marshal(m); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}
func parseInt(s string) (i int) {
	i, _ = strconv.Atoi(s)
	return
}

func parseInt64(s string) (i int64) {
	i, _ = strconv.ParseInt(s, 10, 64)
	return
}

var metricMapping = map[string]map[string]string{
	"riak": map[string]string{
		"VNodeGets":        "VNodeGets",
		"VNodeGetsTotal":   "VNodeGetsTotal",
		"VNodePuts":        "VNodePuts",
		"RingMembersCount": "RingMembersCount",
	},
}

func AllMetricKeys() (ret []string) {
	ret = make([]string, 0)
	for prefix, mapping := range metricMapping {
		for k, _ := range mapping {
			ret = append(ret, prefix+"."+k)
		}
	}
	sort.Strings(ret)
	return
}

func MakeMetric(t, internalKey string, v string, tags map[string]string) (m *Metric) {
	mapping, ok := metricMapping[t]
	if !ok {
		panic("key " + t + " can not defined in metric mapping")
	}
	key, ok := mapping[internalKey]

	if !ok {
		panic("key " + internalKey + " not defined")
	}
	if tags == nil {
		tags = map[string]string{}
	}
	m = &Metric{Key: t + "." + key, Value: v, Tags: tags}
	return
}

func (m *Metric) Ascii(t time.Time, hostname string) (r string) {
	r = fmt.Sprintf("%s %d %s", m.Key, t.Unix(), m.Value)
	m.Tags["host"] = hostname
	for k, v := range m.Tags {
		if len(v) > 0 {
			r = r + " " + k + "=" + v
		}
	}
	return
}

func (m *Metric) NormalizeTag(s string) (r string) {
	re := regexp.MustCompile("(^\\()|(\\)$)")
	re2 := regexp.MustCompile("[^\\w]+")
	return re2.ReplaceAllString(re.ReplaceAllString(s, ""), "_")
}

/*func (m *Metric) Graphite(t time.Time, hostname string) (r string) {
	key := m.Key
	if cpu_id, ok := m.Tags["cpu_id"]; ok {
		key = strings.Replace(key, "cpu", "cpu"+cpu_id, 1)
	}
	if strings.HasPrefix(key, "disk.") {
		if name, ok := m.Tags["name"]; ok {
			key = strings.Replace(key, "disk.", "disk."+name+".", 1)
		}
	}
	if strings.HasPrefix(key, "df.") {
		if name, ok := m.Tags["file_system"]; ok {
			if strings.HasPrefix(name, "/") {
				key = strings.Replace(key, "df.", "df."+strings.Join(strings.Split(name, "/")[1:], ".")+".", 1)
			}
		}
	}
	if strings.HasPrefix(key, "processes.") {
		pid, _ := m.Tags["pid"]
		name, _ := m.Tags["name"]
		if name != "" && pid != "" {
			key = strings.Replace(key, "processes.", "processes."+name+"."+pid+".", 1)
		}
	}
	r = fmt.Sprintf("metrix.hosts.%s.%s %d %d", hostname, key, m.Value, t.Unix())
	return
}
*/
