package api

import (
	//"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/googollee/go-socket.io"
	"github.com/megamsys/libgo/cmd"
	//"github.com/megamsys/vertice/provision"
)

const (
	LOG      = "log"
	VNC      = "vnc"
	NOTFOUND = "Unsupported category"
)

func socketHandler(server *socketio.Server) {
	server.On("connection", func(so socketio.Socket) {
		log.Debugf(cmd.Colorfy("  > [socket] ", "blue", "", "bold") + fmt.Sprintf("Connecting new client : %s", so.Id()))
		so.On("category_connect", func(category string) {
			switch category {
			case LOG:
				logHandler(so)
			case VNC:
				//vncHandler(so)
			default:
				so.Emit("error", NOTFOUND)
			}
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Debugf(cmd.Colorfy("  > [socket] ", "red", "", "bold") + fmt.Sprintf("Error conecting socket server : %v", err))
	})
}
