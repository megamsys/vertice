// Copyright 2015 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package carton

import (
	"time"

	"gopkg.in/check.v1"

)

type WriterSuite struct {
	conn *db.Storage
}

var _ = check.Suite(&WriterSuite{})


func (s *WriterSuite) TestLogWriter(c *check.C) {
	writer := LogWriter{App: &a}
	data := []byte("ble")
	_, err = writer.Write(data)
	c.Assert(err, check.IsNil)
	}


func (s *WriterSuite) TestLogWriterShouldReturnTheDataSize(c *check.C) {
	writer := LogWriter{App: &a}
	data := []byte("ble")
	n, err := writer.Write(data)
	c.Assert(err, check.IsNil)
	c.Assert(n, check.Equals, len(data))
}

func (s *WriterSuite) TestLogWriterAsync(c *check.C) {
	writer := LogWriter{App: &a}
	writer.Async()
	data := []byte("ble")
	_, err = writer.Write(data)
	c.Assert(err, check.IsNil)
	writer.Close()
	err = writer.Wait(5 * time.Second)
	c.Assert(err, check.IsNil)
}


func (s *WriterSuite) TestLogWriterAsyncCopySlice(c *check.C) {
	writer := LogWriter{App: &a}
	writer.Async()
	for i := 0; i < 100; i++ {
		data := []byte("ble")
		_, err = writer.Write(data)
		data[0] = 'X'
		c.Assert(err, check.IsNil)
	}
	writer.Close()
	err = writer.Wait(5 * time.Second)
	c.Assert(err, check.IsNil)
	instance := App{}
	err = s.conn.Apps().Find(bson.M{"name": a.Name}).One(&instance)
	logs, err := instance.LastLogs(100, Applog{})
	c.Assert(err, check.IsNil)
	c.Assert(logs, check.HasLen, 100)
	for i := 0; i < 100; i++ {
		c.Assert(logs[i].Message, check.Equals, "ble")
		c.Assert(logs[i].Source, check.Equals, "tsuru")
	}
}
