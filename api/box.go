package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/megamd/api/context"
	"github.com/megamsys/megamd/provision"
	"golang.org/x/net/websocket"
)

func logs(ws *websocket.Conn) {
	var err error
	defer func() {
		data := map[string]interface{}{}
		if err != nil {
			data["error"] = err.Error()
			log.Error(err.Error())
		} else {
			data["error"] = nil
		}
		msg, _ := json.Marshal(data)
		ws.Write(msg)
		ws.Close()
	}()
	req := ws.Request()
	_ = context.GetAuthToken(req)

	scanner := bufio.NewScanner(ws)
	for scanner.Scan() {
		var entry provision.Boxlog
		data := bytes.TrimSpace(scanner.Bytes())
		if len(data) == 0 {
			continue
		}
		_ = json.Unmarshal(data, &entry)
	}

	err = scanner.Err()
	if err != nil {
		err = fmt.Errorf("wslogs: waiting for log data: %s", err)
		return
	}

	l, _ := provision.NewLogListener(&provision.Box{})
	if err != nil {
		//log the errror.
		return
	}
	LogTracker.add(l)
	defer func() {
		LogTracker.remove(l)
		l.Close()
	}()
	for log := range l.B {
		//wait on the channel, and push it to the ws.
		fmt.Printf("%v", log)
	}

}
