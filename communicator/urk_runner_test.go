package communicator

import (
   "fmt"
    "github.com/megamsys/libgo/cmd"
    "testing"
    //"errors"
  	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

var _ = check.Suite(&HostSuite{})

type HostSuite struct {
}


type Temp struct {
	LvmInstall       bool        `json:"lvminstall"`
	FormatPartitions bool        `json:"formatfartiontions"`
	DeletePartitions bool        `json:"deletepartitions"`
  HostInfo         bool        `json:"hostinfo"`
  NilavuInstall    bool  `json:"nilavuinstall"`
  GatewayInstall   bool  `json:"gatewayinstall"`
  MegamdInstall    bool  `json:"megamdinstall"`
	ParseExistLvm    bool        `json:"parseexistlvm"`
	RemoveLvm        bool        `json:"removelvm"`
	Host             string      `json:"ipaddress"`
	Disks            cmd.MapFlag `json:"osds"`
	Mount            cmd.MapFlag `json:"mount"`
	LvPaths          cmd.MapFlag `json:"lvpaths"`
	VgName           string      `json:"vgname"`
	PvName           string      `json:"pvname"`
	Username         string      `json:"username"`
	Password         string      `json:"password"`
	UserMail         string      `json:"email"`
  Inputs           map[string]string `json:"inputs"`
}

func (s *HostSuite) TestSetConf(c *check.C) {


}


func (s *HostSuite) TestRunners(c *check.C) {
  z := &Temp{
    HostInfo: true,
    Host: "192.168.0.100:22",
    Username: "vijay",
    Password: "speed",
    Inputs: map[string]string{"email": "info@megam.io", "host_id": "123"},
  }
  runner, err := NewUrkRunner(z, z.Inputs)
  if runner !=nil {
		if r, ok := runner.(UrkRunner); ok {
			_, err = r.Run([]string{"HostInfo","NilavuInstall","MegamdInstall"})
	   c.Assert(err, check.IsNil)
     ss := r.OutBuffer.String()
     fmt.Println(ss)
		}
	}


  //  hi := HostInfos{}
  //  hi.ParserCommand(ss)
  //  fmt.Println(hi)
   fmt.Println("Error  :",err)
   //err = errors.New("testing")
	c.Assert(err, check.IsNil)
}
// */
