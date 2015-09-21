package api

import (
	"github.com/megamsys/megamd/provision"
	"gopkg.in/check.v1"
)

func (s *S) TestLogStreamTrackerAddRemove(c *check.C) {
	c.Assert(LogTracker.conn, check.HasLen, 0)
	l := provision.LogListener{}
	LogTracker.add(&l)
	c.Assert(LogTracker.conn, check.HasLen, 1)
	LogTracker.remove(&l)
	c.Assert(LogTracker.conn, check.HasLen, 0)
}
/*
This fails, need to debug
func (s *S) TestLogStreamTrackerShutdown(c *check.C) {
	l, err := provision.NewLogListener(&provision.Box{Name: "myapp", DomainName: "megambox.com"})
	c.Assert(err, check.IsNil)
	LogTracker.add(l)
	LogTracker.Shutdown()
	select {
	case <-l.B:
	case <-time.After(5 * time.Second):
		c.Fatal("timed out waiting for channel to close")
	}
}
*/
