package api

import (
	//	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/googollee/go-socket.io"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/vertice/provision"
)

func logHandler(so socketio.Socket) {
	so.On("logInit", func(msg string) {
		var entry provision.Box
		entry.Name = msg
		l, _ := provision.NewLogListener(&entry)
		go func() {
			so.On("logDisconnect", func(data string) {
				l.Close()
				log.Debugf(cmd.Colorfy("  > [nsqd] unsub   ", "blue", "", "bold") + fmt.Sprintf("Unsubscribing from the Queue"))
			})
			for logbox := range l.B {
				//	logData, _ := json.Marshal(logbox)
				so.Emit(entry.Name, logbox)
			}
		}()

	})

	so.On("disconnection", func() {
		log.Debugf(cmd.Colorfy("  > [socket] ", "blue", "", "bold") + fmt.Sprintf("Disconneted client : %s", so.Id()))
	})

}

/*import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/vertice/provision"
	"net/http"
)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
func logs(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Error in socket connection")
		return err
	}
	messageType, p, err := conn.ReadMessage()
	if err != nil {
		return err
	}
	var entry provision.Box
	_ = json.Unmarshal(p, &entry)
	l, _ := provision.NewLogListener(&entry)
	go func() {
		if _, _, err := conn.NextReader(); err != nil {
			conn.Close()
			l.Close()
			log.Debugf(cmd.Colorfy("  > [nsqd] unsub   ", "blue", "", "bold") + fmt.Sprintf("Unsubscribing from the Queue"))
		}
	}()
	for logbox := range l.B {
		logData, _ := json.Marshal(logbox)
		conn.WriteMessage(messageType, logData)
	}
	return nil
}*/
