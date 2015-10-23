package api

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megamd/provision"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func logs(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Error in socket connection")
		return err
	}
	messageType, p, err := conn.ReadMessage()
	var entry provision.Box

	_ = json.Unmarshal(p, &entry)
	if err != nil {
		return err
	}

	l, _ := provision.NewLogListener(&entry)

	go func() {
		if _, _, err := conn.NextReader(); err != nil {
			conn.Close()
			l.Close()
			log.Debugf(cmd.Colorfy("  > [amqp] unsub   ", "blue", "", "bold") + fmt.Sprintf("Unsubscribing from the Queue"))
		}
	}()

	for log := range l.B {
		logData, _ := json.Marshal(log)
		go conn.WriteMessage(messageType, logData)
	}

	return nil
}
