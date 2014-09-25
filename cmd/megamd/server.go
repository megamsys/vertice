package main

import (
	"fmt"
	"runtime"
	"time"
    "github.com/tsuru/config"
    log "code.google.com/p/log4go"
	"github.com/megamsys/megamd/cmd/megamd/server"
)

func StartDaemon(dry bool) {

	runtime.GOMAXPROCS(runtime.NumCPU())

	version, _ := config.GetString("version")

	log.Info("Starting Megamd Server %s...", version)

	fmt.Printf(`
+------------------------------------------------------------+
|    e   e                                               888 |
|   d8b d8b     ,e e,   e88 888  ,"Y88b 888 888 8e   e88 888 |
|  e Y8b Y8b   d88 88b d888 888 "8" 888 888 888 88b d888 888 |
| d8b Y8b Y8b  888   , Y888 888 ,ee 888 888 888 888 Y888 888 |
|d888b Y8b Y8b  "YeeP"  "88 888 "88 888 888 888 888  "88 888 |
|                        ,  88P                              |
|                       "8",P"                               |
+------------------------------------------------------------+
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
		log.Error("ListenAndServe failed: ", err)
	}
}
