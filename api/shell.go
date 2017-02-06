package api

import (
	"net/http"
	//"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/errors"
	//"github.com/megamsys/vertice/api/context"
	"github.com/megamsys/vertice/carton"
	"github.com/megamsys/vertice/provision"
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
	/*token := context.GetAuthToken(r)
	if token == nil {
		httpErr = &errors.HTTP{
			Code:    http.StatusUnauthorized,
			Message: "no token provided",
		}
		return
	}*/
	/*user, err := token.User()
	if err != nil {
		httpErr = &errors.HTTP{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return
	}*/
	assembly_id := r.URL.Query().Get(":id") //send the assembly_id
	asmsid := r.URL.Query().Get(":asmsid")
	car, err := getBox(asmsid, assembly_id)
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
	boxId := r.URL.Query().Get(":id")
	//width, _ := strconv.Atoi(r.URL.Query().Get("width"))
	//height, _ := strconv.Atoi(r.URL.Query().Get("height"))
	//term := r.URL.Query().Get("term")
	width := 140
	height := 38
	term := "xterm"
	log.Debugf("%s %s %s %s", boxId, width, height, term)

	for _, box := range *car.Boxes {
		opts := provision.ShellOptions{
			Box:    &box,
			Conn:   ws,
			Width:  width,
			Height: height,
			Unit:   boxId,
			Term:   term,
		}
		err = carton.ProvisionerMap[box.Provider].Shell(opts) //BUG: we need get the provisioner of the correct provider
		if err != nil {
			httpErr = &errors.HTTP{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
	}
}


//Return the Box object ?  Get the carton, and make a Box
func getBox(asmsid string, id string) (*carton.Carton, error) {
	c, err := carton.NewCarton(asmsid, id, "rajthilak@megam.io")
	return c, err
}
