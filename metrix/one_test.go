package metrix

import (
	"fmt"

	"gopkg.in/check.v1"
)

// Ensure the configuration can be parsed.
func (s *S) TestParseOpenNebulaCollector(c *check.C) {
	mh := &MetricHandler{}
	on := &OpenNebula{RawStatus: s.testjson}

	all, _ := mh.Collect(on)
	mapped := map[string]*Metric{}
	for _, m := range all {
		//	fmt.Println(m)
		mapped[m.Key] = m
	}
	fmt.Println(mapped)
	c.Assert(mapped["one.accounts_id"].Value, check.Equals, "")
	c.Assert(len(mapped["one.assembly_name"].Value) > 0, check.Equals, true)
	c.Assert(len(mapped["one.assembly_id"].Value) > 0, check.Equals, true)
	c.Assert(len(mapped["one.assemblies_id"].Value) > 0, check.Equals, true)
	c.Assert(len(mapped["one.status"].Value) > 0, check.Equals, true)
	c.Assert(len(mapped["one.system"].Value) > 0, check.Equals, true)
	c.Assert(len(mapped["one.cpu"].Value) > 0, check.Equals, true)
	c.Assert(len(mapped["one.cpu_cost"].Value) > 0, check.Equals, true)
	c.Assert(len(mapped["one.memory"].Value) > 0, check.Equals, true)
	c.Assert(len(mapped["one.memory_cost"].Value) > 0, check.Equals, true)
	c.Assert(len(mapped["one.audit_period_beginning"].Value) > 0, check.Equals, true)
	c.Assert(len(mapped["one.audit_period_ending"].Value) > 0, check.Equals, true)
	c.Assert(len(mapped["one.audit_period_delta"].Value) > 0, check.Equals, true)
	c.Assert(len(mapped) >= 14, check.Equals, true)
}
