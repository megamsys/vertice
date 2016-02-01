package alerts

import (
	"bytes"
	"strings"

	log "github.com/Sirupsen/logrus"
	mailgun "github.com/mailgun/mailgun-go"
)

const (
	CREATE EventAction = iota
	DESTROY
	STATUS
	DEDUCT
	ONBOARD
	RESET
	INVITE
	BALANCE

	MAILGUN = "mailgun"
	SLACK   = "slack"
	INFOBIP = "infobip"

	API_KEY        = "api_key"
	DOMAIN         = "domain"
	TOKEN          = "token"
	CHANNEL        = "channel"
	USERNAME       = "username"
	PASSWORD       = "password"
	APPLICATION_ID = "application_id"
	MESSAGE_ID     = "message_id"
)

type Notifier interface {
	Notify(eva EventAction, m map[string]string) error
	satisfied() bool
}

type EventAction int

func (v *EventAction) String() string {
	switch *v {
	case CREATE:
		return "create"
	case DESTROY:
		return "destroy"
	case STATUS:
		return "status"
	case DEDUCT:
		return "deduct"
	case ONBOARD:
		return "onboard"
	case RESET:
		return "reset"
	case INVITE:
		return "invite"
	case BALANCE:
		return "bal"
	default:
		return "arrgh"
	}
}

type mailgunner struct {
	api_key string
	domain  string
}

func NewMailgun(m map[string]string) Notifier {
	return &mailgunner{
		api_key: m[API_KEY],
		domain:  m[DOMAIN],
	}
}

func (m *mailgunner) satisfied() bool {
	return true
}

func (m *mailgunner) Notify(eva EventAction, mp map[string]string) error {
	if !m.satisfied() {
		return nil
	}
	var w bytes.Buffer
	var sub string
	var err error

	switch eva {
	case ONBOARD:
		err = onboardTemplate.Execute(&w, mp)
		sub = "Ahoy. Welcome aboard!"
	case RESET:
		err = resetTemplate.Execute(&w, mp)
		sub = "You have fat finger.!"
	case INVITE:
		err = inviteTemplate.Execute(&w, mp)
		sub = "Lets party!"
	case BALANCE:
		err = balanceTemplate.Execute(&w, mp)
		sub = "Piggy bank!"
	case CREATE:
		err = onboardTemplate.Execute(&w, mp)
		sub = "Up!"
	case DESTROY:
		err = onboardTemplate.Execute(&w, mp)
		sub = "Nuked"
	default:
		break
	}
	if err != nil {
		return err
	}
	m.Send(w.String(), "", sub, mp["email"])
	return nil
}

func (m *mailgunner) Send(msg string, sender string, subject string, to string) error {
	if len(strings.TrimSpace(sender)) <= 0 {
		sender = "Kishore CEO <nkishore@megam.io>"
	}
	mg := mailgun.NewMailgun(m.domain, m.api_key, "")
	g := mailgun.NewMessage(
		sender,
		subject,
		"You are in !",
		to,
	)
	g.SetHtml(msg)
	g.SetTracking(true)
	_, id, err := mg.Send(g)
	if err != nil {
		return err
	}
	log.Infof("Mailgun sent %s", id)
	return nil
}

/*return map[string]interface{}{
		"email":     "nkishore@megam.io",
		"logo":      "vertice.png",
		"nilavu":    "console.megam.io",
		"click_url": "https://console.megam.io/reset?email=test@email.com&resettoken=9090909090",
		"days":      "20",
		"cost":      "$12",
}*/
