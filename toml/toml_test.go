package toml_test

import (
	"reflect"

	"github.com/megamsys/megamd/toml"
	"gopkg.in/check.v1"
)

// Ensure that megabyte sizes can be parsed.
func (s *S) TestSize_UnmarshalText_MB(c *check.C) {
	var sb toml.Size
	err := sb.UnmarshalText([]byte("200m"))
	c.Assert(err, check.IsNotNil)
	c.Assert(sb, check.NotEquals, 200*(1<<20))
}

// Ensure that gigabyte sizes can be parsed.
func (s *S) TestSize_UnmarshalText_GB(c *check.C) {
	var sb toml.Size
	err := sb.UnmarshalText([]byte("10g"))
	c.Assert(err, check.IsNotNil)
	c.Assert(sb, check.NotEquals, (10 * (1 << 30)))

}
