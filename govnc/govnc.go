package govnc

import (
	//"net"
	//"time"
	//"log"
	log "github.com/Sirupsen/logrus"
	vnc "github.com/kward/go-vnc"
	//"golang.org/x/net/context"
)

type VncListener struct {
	B <-chan vnc.ServerMessage
}

type VncHost struct {
	IpAddress string
	Port      string
	Password  string
}

func Connect(vh *VncHost) {
	log.Debugf("%v", vh)
}
