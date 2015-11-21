package api

import (
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/errors"
	"github.com/megamsys/megamd/api/context"
	"github.com/megamsys/megamd/carton"
	"github.com/megamsys/megamd/provision"
	"golang.org/x/net/websocket"
)

func remoteShellHandler(ws *websocket.Conn) {
	var httpErr *errors.HTTP
	defer func() {
		defer ws.Close()
		if httpErr != nil {
			var msg string
			switch httpErr.Code {
			case http.StatusUnauthorized:
				msg = "no token provided or session expired, please login again\n"
			default:
				msg = httpErr.Message + "\n"
			}
			ws.Write([]byte("Error: " + msg))
		}
	}()
	r := ws.Request()
	token := context.GetAuthToken(r)
	if token == nil {
		httpErr = &errors.HTTP{
			Code:    http.StatusUnauthorized,
			Message: "no token provided",
		}
		return
	}
	user, err := token.User()
	if err != nil {
		httpErr = &errors.HTTP{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}
	assembly_id := r.URL.Query().Get(":id") //send the assembly_id
	box, err := getBox(assembly_id)
	if err != nil {
		if herr, ok := err.(*errors.HTTP); ok {
			httpErr = herr
		} else {
			httpErr = &errors.HTTP{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		return
	}
	boxId := r.URL.Query().Get("id")
	width, _ := strconv.Atoi(r.URL.Query().Get("width"))
	height, _ := strconv.Atoi(r.URL.Query().Get("height"))
	term := r.URL.Query().Get("term")
	log.Debugf("%s %s %s %s %s", user, boxId, width, height, term)

	opts := provision.ShellOptions{
		Box:    box,
		Conn:   ws,
		Width:  width,
		Height: height,
		Unit:   boxId,
		Term:   term,
	}
	err = carton.Provisioner.Shell(opts) //BUG: we need get the provisioner of the correct provider
	if err != nil {
		httpErr = &errors.HTTP{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
}

//Return the Box object ?  Get the carton, and make a Box
func getBox(a string) (*provision.Box, error) {
	return nil, nil
}
