package api

import (
	"encoding/json"
	"fmt"
	//"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/googollee/go-socket.io"
	"github.com/megamsys/vertice/provision"
	"github.com/megamsys/libgo/cmd"
	//"encoding/gob"
	//"net/http"
)

func logHandler(so socketio.Socket) {
  fmt.Println("---------------------")
	so.On("logConnect", func(name string) {
		fmt.Println("----------connect-----------")
		fmt.Println(name)
			//so.Join(so.Id())
	})

	so.On("logInit", func(msg string) {
		fmt.Println("----------init-----------")
		fmt.Println(msg)
		var entry provision.Box
		//var buf bytes.Buffer
    //enc := gob.NewEncoder(&buf)
    //err := enc.Encode(msg)
    //if err != nil {
				//so.Emit("error", err)
        //return
    //}
		fmt.Println("---------parse------------")
		//_ = json.Unmarshal([]byte(msg), &entry)
		entry.Name = msg
		l, _ := provision.NewLogListener(&entry)

		fmt.Println("-------end--------------")
		fmt.Println(l.B)
		fmt.Println(entry.Name)
		//go logclose(so, l)

		go func() {
			so.On("logDisconnect", func(data string) {
				fmt.Println("djbvjfbvkdfjbvkdfjvbdfkjvbkjkn")
				l.Close()
				log.Debugf(cmd.Colorfy("  > [nsqd] unsub   ", "blue", "", "bold") + fmt.Sprintf("Unsubscribing from the Queue"))
			})
		}()

		for logbox := range l.B {
			fmt.Println(logbox)
			logData, err := json.Marshal(logbox)
			fmt.Println(err)
			fmt.Println(logData)
			so.Emit(entry.Name, logData)
		}

	})
}

/*func logclose(so socketio.Socket, l *provision.LogListener) {
	fmt.Println("-----------disconnect enter----------")
	so.On("logDisconnect", func(data string) {
		fmt.Println(data)
			so.Leave(so.Id())
			l.Close()
			log.Debugf(cmd.Colorfy("  > [nsqd] unsub   ", "blue", "", "bold") + fmt.Sprintf("Unsubscribing from the Queue"))
			return
	})
}*/


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
