package metrix

import (
	"io/ioutil"
	"testing"
	//"fmt"

	"github.com/BurntSushi/toml"
	"github.com/megamsys/vertice/meta"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type S struct {
	testjson []byte
	testxml  []byte
	cm       meta.Config
	sensor   *Sensor
}

var _ = check.Suite(&S{})

func (s *S) SetUpSuite(c *check.C) {
	b, e := ioutil.ReadFile("fixtures/one.xml")
	c.Assert(e, check.IsNil)
	s.testxml = b
	c.Assert(s.testxml, check.NotNil)

	sen := &Sensor{
		Id:                   "22a0fd6ac83411e6aa700cc47a143ea0",
		AccountId:            "rajthilak@megam.io",
		SensorType:           "compute.instance.exists",
		AssemblyId:           "ASM8963235353929524317",
		AssemblyName:         "lingering-feather-2706.megambox.com",
		AssembliesId:         "AMS5473376480118161453",
		Node:                 "146.0.247.2",
		System:               "one",
		Status:               "Active",
		Source:               "one",
		Message:              "vm billing",
		AuditPeriodBeginning: "1970-01-01 01:00:00 +0100 CET",
		AuditPeriodEnding:    "1970-01-01 01:10:00 +0100 CET",
		AuditPeriodDelta:     "4.117788035058207E+05",
	}
	s.sensor = sen

	var cm meta.Config
	if _, err := toml.Decode(`
		api = "http://192.168.1.4:9000/v2"
		master_key = "3b8eb672aa7c8db82e5d34a0744740b20ed59e1f6814cfb63364040b0994ee3f"
		nsqd = ["192.168.1.100:4150"]
	`, &cm); err != nil {
		c.Fatal(err)
	}
	cm.MkGlobal()
	s.cm = cm
}
