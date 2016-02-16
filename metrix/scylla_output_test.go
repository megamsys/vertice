package metrix

import (
	"gopkg.in/check.v1"
)

func (s *S) TestStoreCollectedr(c *check.C) {
	c.Skip("Fix: Ping riak, and then decide to skip sensors")

	mh := &MetricHandler{}
	on := &OpenNebula{RawStatus: s.testxml}
	all, _ := mh.Collect(on)
	c.Assert(all, check.NotNil)

	o := OutputHandler{
		RiakAddress: s.cm.Riak,
	}
	c.Assert(o, check.NotNil)
	err := o.WriteMetrics(all)
	c.Assert(err, check.IsNil)
}
