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
	"log"
	"os"
)

func NewMultiLogger(loggers ...Logger) Logger {
	return &multiLogger{loggers: loggers}
}

type multiLogger struct {
	loggers []Logger
}

func (m *multiLogger) Debug(message string) {
	for _, logger := range m.loggers {
		logger.Debug(message)
	}
}

func (m *multiLogger) Error(message string) {
	for _, logger := range m.loggers {
		logger.Error(message)
	}
}

func (m *multiLogger) Fatal(message string) {
	for _, logger := range m.loggers {
		logger.Error(message)
	}
	os.Exit(1)
}

func (m *multiLogger) Debugf(format string, v ...interface{}) {
	for _, logger := range m.loggers {
		logger.Debugf(format, v...)
	}
}

func (m *multiLogger) Errorf(format string, v ...interface{}) {
	for _, logger := range m.loggers {
		logger.Errorf(format, v...)
	}
}

func (m *multiLogger) Fatalf(format string, v ...interface{}) {
	for _, logger := range m.loggers {
		logger.Errorf(format, v...)
	}
	os.Exit(1)
}

func (m *multiLogger) GetStdLogger() *log.Logger {
	return m.loggers[0].GetStdLogger()
}
