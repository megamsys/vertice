package eventsd

import (
	b64 "encoding/base64"
	"github.com/megamsys/libgo/events"
	"github.com/megamsys/libgo/events/alerts"
	constants "github.com/megamsys/libgo/utils"
	"github.com/megamsys/opennebula-go/api"
	"github.com/megamsys/opennebula-go/users"
	"github.com/megamsys/vertice/subd/deployd"
)

type Handler struct {
	Deployd      *deployd.Config
	d            *Config
	EventChannel chan bool
}

func NewHandler(c *Config) *Handler {
	return &Handler{d: c}

}

func (h *Handler) serveNSQ(e *events.Event, email string) error {
	if h.isOnboard(e) {
		e.EventData.M[constants.NILAVU_PASSWORD] = h.decryptBase64(e.EventData.M[constants.NILAVU_PASSWORD])
	}
	if err := events.W.Write(e); err != nil {
		return err
	}
	return nil
}

func (m *Handler) isOnboard(e *events.Event) bool {
	return e.EventAction == alerts.ONBOARD && e.EventType == constants.EventUser
}

func (m *Handler) decryptBase64(pswd string) string {
	pwd, _ := b64.StdEncoding.DecodeString(pswd)
	return string(pwd)
}

func (h *Handler) userProviderOnboard(email string) error {
	for _, region := range h.Deployd.One.Regions {
		cm := make(map[string]string)
		cm[api.ENDPOINT] = region.OneEndPoint
		cm[api.USERID] = region.OneUserid
		cm[api.PASSWORD] = region.OnePassword
		client, _ := api.NewClient(cm)
		u := users.User{
			UserName:   email,
			Password:   region.OneMasterKey,
			AuthDriver: "core",
			GroupIds:   []int{0},
		}
		vm := users.UserTemplate{
			T:     client,
			Users: u,
		}
		_, err := vm.CreateUsers()
		if err != nil {
			return err
		}
	}
	return nil
}
