package govnc

import (
	//"net"
	//"time"
	//"log"
	vnc "github.com/kward/go-vnc"
	log "github.com/Sirupsen/logrus"
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
	log.Debugf("%v",vh)
}
