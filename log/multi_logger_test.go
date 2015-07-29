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

	"gopkg.in/check.v1"
)

type MultiLoggerSuite struct {
	logger     Logger
	buf1, buf2 bytes.Buffer
}

var _ = check.Suite(&MultiLoggerSuite{})

func (s *MultiLoggerSuite) SetUpTest(c *check.C) {
	s.logger = NewMultiLogger(
		NewWriterLogger(&s.buf1, true),
		NewWriterLogger(&s.buf2, true),
	)
}

func (s *MultiLoggerSuite) TearDownTest(c *check.C) {
	s.buf1.Reset()
	s.buf2.Reset()
}

func (s *MultiLoggerSuite) TestDebug(c *check.C) {
	s.logger.Debug("something went wrong")
	c.Check(s.buf1.String(), check.Matches, `(?m)^.*DEBUG: something went wrong$`)
	c.Check(s.buf1.String(), check.Matches, `(?m)^.*DEBUG: something went wrong$`)
}

func (s *MultiLoggerSuite) TestDebugf(c *check.C) {
	s.logger.Debugf("something went wrong: %q", "this")
	c.Check(s.buf1.String(), check.Matches, `(?m)^.*DEBUG: something went wrong: "this"$`)
	c.Check(s.buf1.String(), check.Matches, `(?m)^.*DEBUG: something went wrong: "this"$`)
}

func (s *MultiLoggerSuite) TestError(c *check.C) {
	s.logger.Error("something went wrong")
	c.Check(s.buf1.String(), check.Matches, `(?m)^.*ERROR: something went wrong$`)
	c.Check(s.buf1.String(), check.Matches, `(?m)^.*ERROR: something went wrong$`)
}

func (s *MultiLoggerSuite) TestErrorf(c *check.C) {
	s.logger.Errorf("something went wrong: %q", "this")
	c.Check(s.buf1.String(), check.Matches, `(?m)^.*ERROR: something went wrong: "this"$`)
	c.Check(s.buf1.String(), check.Matches, `(?m)^.*ERROR: something went wrong: "this"$`)
}
