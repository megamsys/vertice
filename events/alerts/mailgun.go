package alerts

import (
	"bytes"
	"strings"

	log "github.com/Sirupsen/logrus"
	mailgun "github.com/mailgun/mailgun-go"
)

const (
	ONBOARD int = iota
	RESET
	INVITE
	BALANCE
	MACHINE_CREATED
	MACHINE_DESTROYED
)

type Notifier interface {
	Notify(evt string) error
	satisfied() bool
}

type mailgunner struct {
	api_key string
	domain  string
}

func NewMailgun(ak string, do string) Notifier {
	return &mailgunner{
		api_key: ak,
		domain:  do,
	}
}

func (m *mailgunner) satisfied() bool {
	return true
}

func (m *mailgunner) Notify(evt string) error {
	if !m.satisfied() {
		return nil
	}
	var w bytes.Buffer
	var sub string
	var err error

	eventAction := 1

	switch eventAction {
	case ONBOARD:
		err = onboardTemplate.Execute(&w, onboard())
		sub = "Ahoy. Welcome aboard!"
	case RESET:
		err = resetTemplate.Execute(&w, reset())
		sub = "You have fat finger.!"
	case INVITE:
		err = inviteTemplate.Execute(&w, fake())
		sub = "Lets party!"
	case BALANCE:
		err = balanceTemplate.Execute(&w, fake())
		sub = "Piggy bank!"
	case MACHINE_CREATED:
		err = onboardTemplate.Execute(&w, fake())
		sub = "Up!"
	case MACHINE_DESTROYED:
		err = onboardTemplate.Execute(&w, fake())
		sub = "Nuked"
	default:
		break
	}
	if err != nil {
		return err
	}
	m.Send(w.String(), "", sub, "info@megam.io")
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

func onboard() map[string]interface{} {
	return map[string]interface{}{
		"email":     "nkishore@megam.io",
		"logo":      "vertice.png",
		"nilavu":    "console.megam.io",
		"click_url": "https://console.megam.io/reset?email=test@email.com&resettoken=9090909090",
		"days":      "20",
		"cost":      "$12",
	}
}

func reset() map[string]interface{} {
	return map[string]interface{}{
		"email":     "nkishore@megam.io",
		"logo":      "vertice.png",
		"nilavu":    "console.megam.io",
		"click_url": "https://console.megam.io/reset?email=test@email.com&resettoken=9090909090",
		"days":      "20",
		"cost":      "$12",
	}
}

func fake() map[string]interface{} {
	return map[string]interface{}{
		"email":     "nkishore@megam.io",
		"logo":      "vertice.png",
		"nilavu":    "console.megam.io",
		"click_url": "https://console.megam.io/reset?email=test@email.com&resettoken=9090909090",
		"days":      "20",
		"cost":      "$12",
	}
}
