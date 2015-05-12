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
	"github.com/megamsys/megamd/cmd/megamd/server"
	"github.com/megamsys/megamd/global"
	"github.com/tsuru/config"
	"runtime"
	"time"
)

func StartDaemon(dry bool) {

	runtime.GOMAXPROCS(runtime.NumCPU())

	version, _ := config.GetString("version")
    global.LogFormatter()
	global.LOG.Info("Starting Megamd Server %s...", version)		

	global.LOG.Notice(`
███╗   ███╗███████╗ ██████╗  █████╗ ███╗   ███╗██████╗ 
████╗ ████║██╔════╝██╔════╝ ██╔══██╗████╗ ████║██╔══██╗
██╔████╔██║█████╗  ██║  ███╗███████║██╔████╔██║██║  ██║
██║╚██╔╝██║██╔══╝  ██║   ██║██╔══██║██║╚██╔╝██║██║  ██║
██║ ╚═╝ ██║███████╗╚██████╔╝██║  ██║██║ ╚═╝ ██║██████╔╝
╚═╝     ╚═╝╚══════╝ ╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝╚═════╝ 
`)

	server, err := server.NewServer()
	if err != nil {
		// sleep for the log to flush
		time.Sleep(time.Second)
		panic(err)
	}

	if err := startProfiler(server); err != nil {
		panic(err)
	}

	err = server.ListenAndServe()
	if err != nil {
		global.LOG.Error("ListenAndServe failed: ", err)
	}
}
