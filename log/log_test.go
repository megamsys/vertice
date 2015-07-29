/*
** Copyright [2013-2015] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */


package log

import (
	"bytes"
	"log"
	"testing"

	"github.com/tsuru/config"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct{}

var _ = check.Suite(&S{})

func newFakeLogger() *bytes.Buffer {
	l := NewFileLogger("/dev/null", true)
	fl, _ := l.(*fileLogger)
	b := &bytes.Buffer{}
	fl.logger = log.New(b, "", 0)
	SetLogger(l)
	return b
}

func (s *S) TestLogError(c *check.C) {
	buf := newFakeLogger()
	defer buf.Reset()
	Error("log anything")
	c.Assert(buf.String(), check.Equals, "ERROR: log anything\n")
}

func (s *S) TestLogErrorf(c *check.C) {
	buf := newFakeLogger()
	defer buf.Reset()
	Errorf("log anything %d", 1)
	c.Assert(buf.String(), check.Equals, "ERROR: log anything 1\n")
}

func (s *S) TestLogErrorWithoutTarget(c *check.C) {
	_ = newFakeLogger()
	defer func() {
		c.Assert(recover(), check.IsNil)
	}()
	Error("log anything")
}

func (s *S) TestLogErrorfWithoutTarget(c *check.C) {
	_ = newFakeLogger()
	defer func() {
		c.Assert(recover(), check.IsNil)
	}()
	Errorf("log anything %d", 1)
}

func (s *S) TestLogDebug(c *check.C) {
	buf := newFakeLogger()
	defer buf.Reset()
	Debug("log anything")
	c.Assert(buf.String(), check.Equals, "DEBUG: log anything\n")
}

func (s *S) TestLogDebugf(c *check.C) {
	buf := newFakeLogger()
	defer buf.Reset()
	Debugf("log anything %d", 1)
	c.Assert(buf.String(), check.Equals, "DEBUG: log anything 1\n")
}

func (s *S) TestWrite(c *check.C) {
	w := &bytes.Buffer{}
	err := Write(w, []byte("teeest"))
	c.Assert(err, check.IsNil)
	c.Assert(w.String(), check.Equals, "teeest")
}

func (s *S) TestInitWithWrongConf(c *check.C) {
	configFile := "testdata/wrongconfig.yml"
	err := config.ReadConfigFile(configFile)
	c.Assert(err, check.IsNil)
	c.Assert(Init, check.PanicMatches, "Your conf is wrong: please see http://docs.megam.io")
}
