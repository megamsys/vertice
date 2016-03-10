package metrix

import (
	"io/ioutil"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/megamsys/vertice/meta"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type S struct {
	testxml []byte
	cm      meta.Config
}

var _ = check.Suite(&S{})

func (s *S) SetUpSuite(c *check.C) {
	b, e := ioutil.ReadFile("fixtures/one.xml")
	c.Assert(e, check.IsNil)
	s.testxml = b
	c.Assert(s.testxml, check.NotNil)

	var cm meta.Config
	if _, err := toml.Decode(`
	scylla = ["localhost:8087"]
	`, &cm); err != nil {
		c.Fatal(err)
	}
	cm.MkGlobal()
	s.cm = cm
}
