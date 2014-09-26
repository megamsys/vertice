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
		log.Error("ListenAndServe failed: ", err)
	}
}
