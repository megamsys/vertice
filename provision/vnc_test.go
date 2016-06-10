package provision

import (
//	"sync"
//	"time"
"fmt"
	"gopkg.in/check.v1"
)
type S struct{}


func (s *S) TestNewLogListener(c *check.C) {
	box := Box{Name: "myapp"}
	l, err := NewVNCListener(&box)
  fmt.Println("*****************")
  fmt.Println(l)
//	defer l.Close()
	c.Assert(err, check.IsNil)
	//c.Assert(l.q, check.NotNil)
	//c.Assert(l.C, check.NotNil)
	///notify("mybox", []interface{}{Boxlog{Message: "123"}})

}
