package metrix

import (
	"gopkg.in/check.v1"
)

// Ensure the configuration can be parsed.
func (s *S) TestParseOpenNebulaCollector(c *check.C) {
	mh := &MetricHandler{}
	on := &OpenNebula{RawStatus: s.testxml}
	all, _ := mh.Collect(on)

	for _, m := range all {
		c.Assert(len(m.Id) > 0, check.Equals, false)
		c.Assert(len(m.AccountId) > 0, check.Equals, true)
		c.Assert(len(m.SensorType) > 0, check.Equals, true)
		c.Assert(len(m.AssemblyId) > 0, check.Equals, true)
		c.Assert(len(m.AssembliesId) > 0, check.Equals, true)
		c.Assert(len(m.AssemblyName) > 0, check.Equals, true)
		c.Assert(len(m.System) > 0, check.Equals, true)
		c.Assert(len(m.Status) > 0, check.Equals, true)
		c.Assert(len(m.Source) > 0, check.Equals, true)
		c.Assert(len(m.Message) > 0, check.Equals, true)
		c.Assert(len(m.AuditPeriodBeginning) > 0, check.Equals, true)
		c.Assert(len(m.AuditPeriodEnding) > 0, check.Equals, true)
		c.Assert(len(m.AuditPeriodDelta) > 0, check.Equals, true)

		for _, me := range m.Metrics {
			c.Assert(len(me.MetricName) > 0, check.Equals, true)
			c.Assert(len(me.MetricValue) > 0, check.Equals, true)
			c.Assert(len(me.MetricUnits) > 0, check.Equals, true)
			c.Assert(len(me.MetricType) > 0, check.Equals, true)

		}

	}
}
