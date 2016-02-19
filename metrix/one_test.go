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
		c.Assert(len(m.AccountsId) > 0, check.Equals, true)
		c.Assert(len(m.Type) > 0, check.Equals, true)
		c.Assert(len(m.CreatedAt) > 0, check.Equals, true)
		c.Assert(len(m.Payload.AssemblyId) > 0, check.Equals, true)
		c.Assert(len(m.Payload.AssembliesId) > 0, check.Equals, true)
		c.Assert(len(m.Payload.AssemblyName) > 0, check.Equals, true)
		c.Assert(len(m.Payload.System) > 0, check.Equals, true)
		c.Assert(len(m.Payload.Status) > 0, check.Equals, true)
		c.Assert(len(m.Payload.Source) > 0, check.Equals, true)
		c.Assert(len(m.Payload.Message) > 0, check.Equals, true)
		c.Assert(len(m.Payload.BeginAudit) > 0, check.Equals, true)
		c.Assert(len(m.Payload.EndAudit) > 0, check.Equals, true)
		c.Assert(len(m.Payload.DeltaAudit) > 0, check.Equals, true)

		for _, me := range m.Payload.Metrics {
			c.Assert(len(me.Key) > 0, check.Equals, true)
			c.Assert(len(me.Value) > 0, check.Equals, true)
			c.Assert(len(me.Units) > 0, check.Equals, true)
			c.Assert(len(me.Type) > 0, check.Equals, true)

		}

	}
}
