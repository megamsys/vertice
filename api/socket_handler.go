package api

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/googollee/go-socket.io"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/vertice/provision"
)

const (
	LOG      = "log"
	VNC      = "vnc"
	NOTFOUND = "Unsupported category"
)

func socketHandler(server *socketio.Server) {
	server.On("connection", func(so socketio.Socket) {
		log.Debugf(cmd.Colorfy("  > [socket] ", "blue", "", "bold") + fmt.Sprintf("Connecting new client : %s", so.Id()))
		//so.Join(so.Id())

		/*so.On("category_connect", func(category string) {
		      fmt.Println(category)
					switch category {
					case LOG:
		        fmt.Println(category)
						go logHandler(so)
					case VNC:
						//vncHandler(so)
					default:
		        so.Emit("error", NOTFOUND)
					}
				})*/

		so.On(LOG, func(name string) {
			fmt.Println(name)
			var entry provision.Box
			entry.Name = name
			l, _ := provision.NewLogListener(&entry)
			go func() {
				so.On("logDisconnect", func(data string) {
					fmt.Println("djbvjfbvkdfjbvkdfjvbdfkjvbkjkn")
					l.Close()
					log.Debugf(cmd.Colorfy("  > [nsqd] unsub   ", "blue", "", "bold") + fmt.Sprintf("Unsubscribing from the Queue"))
				})
			}()
			go func() {
				for logbox := range l.B {
					logData, _ := json.Marshal(logbox)
					fmt.Println("datatatsvchvbkvkkjkj")
					so.Emit(entry.Name, logData)
				}
			}()
		})

		so.On("disconnection", func() {
			fmt.Println("on disconnect")
		})

	})
	server.On("error", func(so socketio.Socket, err error) {
		fmt.Println("error:", err)
	})
}
