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
	testjson []byte
	cm       meta.Config
}

var _ = check.Suite(&S{})

func (s *S) SetUpSuite(c *check.C) {
	b, e := ioutil.ReadFile("fixtures/one.json")
	c.Assert(e, check.IsNil)
	s.testjson = b
	c.Assert(s.testjson, check.NotNil)

	var cm meta.Config
	if _, err := toml.Decode(`
	riak = ["localhost:8087"]
	`, &cm); err != nil {
		c.Fatal(err)
	}
	cm.MkGlobal()
	s.cm = cm
}
