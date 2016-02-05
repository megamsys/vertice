package toml_test

import (
	//	"reflect"
	"testing"

	"github.com/megamsys/vertice/toml"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type S struct{}

var _ = check.Suite(&S{})

// Ensure that megabyte sizes can be parsed.
func (s *S) TestSize_UnmarshalText_MB(c *check.C) {
	var sb toml.Size
	err := sb.UnmarshalText([]byte("200m"))
	c.Assert(err, check.IsNil)
	c.Assert(sb, check.Not(check.Equals), 200*(1<<20))
}

//Ensure that gigabyte sizes can be parsed.
func (s *S) TestSize_UnmarshalText_GB(c *check.C) {
	var sb toml.Size
	err := sb.UnmarshalText([]byte("10g"))
	c.Assert(err, check.IsNil)
	f := int64(10 * (1 << 30))
	c.Assert(sb, check.Not(check.Equals), f)
}
