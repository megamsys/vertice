package api

import (
	//"encoding/json"
	"net/http"
  "fmt"
  log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	//"github.com/megamsys/libgo/cmd"
	//"github.com/googollee/go-socket.io"
	//"github.com/megamsys/vertice/govnc"
)

/*func vnc(w http.ResponseWriter, r *http.Request) error {

	return nil

}*/

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func vnc(w http.ResponseWriter, r *http.Request) error {
  conn, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Errorf("Error in socket connection")
    return err
  }
  s := conn.Subprotocol()
	fmt.Println(s)
/*  messageType, p, err := conn.ReadMessage()

  if err != nil {
    return err
  }

  var entry govnc.VncHost

  _ = json.Unmarshal(p, &entry)
	govnc.Connect(&entry)
  l, _ := govnc.Connect(&entry)

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
conn.WriteMessage(messageType, []byte("hai"))*/
  return nil
}
