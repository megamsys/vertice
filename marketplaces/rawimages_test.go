package marketplaces

import (
	"fmt"
	"gopkg.in/check.v1"
	//  "encoding/json"
)

// func (s *S) TestGetRawimages (c *check.C) {
//   r := new(RawImages)
//   r.AccountId = "vino@gmail.co"
//   r.Id = "RAW5309653640173297515"
//   res, err := r.Get()
//   // repo := &Repos{}
// 	// err = json.Unmarshal([]byte(res.Repository), repo)
//
//   fmt.Println("result :",res)
//   fmt.Println("error :",err)
//   c.Assert(nil, check.NotNil)
// }

func (s *S) TestUpdateRawimages(c *check.C) {
	r := new(RawImages)
	r.AccountId = "vino@gmail.co"
	r.Id = "RAW9058037720298113796" //"RAW5931914238596344628"
	res, err := r.Get()
	fmt.Println("result :", res)
	fmt.Println("error :", err)
	c.Assert(err, check.IsNil)
	res.Status = "TestRawImages"
	err = res.update()
	c.Assert(err, check.IsNil)
	fmt.Println("Post Success")
}
