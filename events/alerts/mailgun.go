package alerts

import (
	log "github.com/Sirupsen/logrus"
	mailgun "github.com/mailgun/mailgun-go"
	"strings"
)

const (
	LAUNCHED EventAction = iota
	DESTROYED
	STATUS
	DEDUCT
	ONBOARD
	RESET
	INVITE
	BALANCE
	INVOICE
	TRANSACTION

	//keys for watchers
	MAILGUN = "mailgun"
	SLACK   = "slack"
	INFOBIP = "infobip"

	//config keys by watchers
	TOKEN          = "token"
	CHANNEL        = "channel"
	USERNAME       = "username"
	PASSWORD       = "password"
	APPLICATION_ID = "application_id"
	MESSAGE_ID     = "message_id"
	API_KEY        = "api_key"
	DOMAIN         = "domain"
	PIGGYBANKS     = "piggybanks"

	//args for notification
	NILAVU    = "nilavu"
	LOGO      = "logo"
	NAME      = "name"
	VERTNAME  = "appname"
	TEAM      = "team"
	VERTTYPE  = "type"
	EMAIL     = "email"
	DAYS      = "days"
	COST      = "cost"
	STARTTIME = "starttime"
	ENDTIME   = "endtime"
)

type Notifier interface {
	Notify(eva EventAction, m map[string]string) error
	satisfied() bool
}

type EventAction int

func (v *EventAction) String() string {
	switch *v {
	case LAUNCHED:
		return "launched"
	case DESTROYED:
		return "destroyed"
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
		return "balance"
	default:
		return "arrgh"
	}
}

type mailgunner struct {
	api_key string
	domain  string
	nilavu  string
	logo    string
}

func NewMailgun(m map[string]string) Notifier {
	return &mailgunner{
		api_key: m[API_KEY],
		domain:  m[DOMAIN],
		nilavu:  m[NILAVU],
		logo:    m[LOGO],
	}
}

func (m *mailgunner) satisfied() bool {
	return true
}

/*{
		"email":     "nkishore@megam.io",
		"logo":      "vertice.png",
		"nilavu":    "console.megam.io",
		"appname": "vertice.megambox.com"
		"type": "torpedo"
		"token": "9090909090",
		"days":      "20",
		"cost":      "$12",
}*/
func (m *mailgunner) Notify(eva EventAction, mp map[string]string) error {
	if !m.satisfied() {
		return nil
	}
	mp[NILAVU] = m.nilavu
	mp[LOGO] = m.logo

	bdy, err := body(eva.String(), mp)
	if err != nil {
		return err
	}
	m.Send(bdy, "", subject(eva), mp[EMAIL])
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

func subject(eva EventAction) string {
	var sub string
	switch eva {
	case ONBOARD:
		sub = "Ahoy. Welcome aboard!"
	case RESET:
		sub = "You have fat finger.!"
	case INVITE:
		sub = "Lets party!"
	case BALANCE:
		sub = "Piggy bank!"
	case LAUNCHED:
		sub = "Up!"
	case DESTROYED:
		sub = "Nuked"
	default:
		break
	}
	return sub
}
