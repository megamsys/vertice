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
package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megamd/cmd/megamd/run"
)

// These variables are populated via the Go linker.
var (
	version string = "0.9.2"
	date    string
	commit  string
)

func init() {
	log.SetOutput(os.Stdout)
}

// Only log debug level when the -v flag is passed.
func cmdRegistry(name string) *cmd.Manager {
	m := cmd.BuildBaseManager(name, version+"  "+date, nil, func(modelvl int) {
		if modelvl >= 1 {
			log.SetLevel(log.DebugLevel)
		}
	})
	m.Register(&run.Start{})
	return m
}

//Run the commands from cli.
func main() {
	name := cmd.ExtractProgramName(os.Args[0])
	manager := cmdRegistry(name)
	manager.Run(os.Args[1:])
}
