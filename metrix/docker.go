package metrix

import (
	//"encoding/xml"
	//"github.com/megamsys/opennebula-go/metrics"
	//"github.com/megamsys/vertice/carton"
	//"io/ioutil"
//	"time"
)

const DOCKER = "docker"

type Swarm struct {
	Url       string
	RawStatus []byte
}

func (s *Swarm) Prefix() string {
	return "docker"
}

func (s *Swarm) Collect(c *MetricsCollection) (e error) {
	 e = s.ReadStatus()
	if e != nil {
		return
	}
  //
	// s, e := s.ParseStatus(b)
	// if e != nil {
	// 	return
	// }
	// on.CollectMetricsFromStats(c, s)
	 return
}

func (s *Swarm) ReadStatus() (e error) {

  return
}
