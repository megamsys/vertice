package events

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/megamd/events/alerts"
	"github.com/megamsys/megamd/subd/eventsd"
)

const (
	MAILGUN = "mailgun"
	SLACK   = "slack"
	INFOBIP = "infobip"
)

var notifiers map[string]alerts.Notifier

type User struct {
	stop chan struct{}
}

func NewUser(e *eventsd.Config) *User {
	register(e)
	return &User{}
}

func register(e *eventsd.Config) {
	notifiers = make(map[string]alerts.Notifier)
	notifiers[MAILGUN] = newMailgun(&e.Mailgun)
	notifiers[INFOBIP] = newInfobip(&e.Infobip)
	notifiers[SLACK] = newSlack(&e.Slack)
}

func newMailgun(m *eventsd.Mailgun) alerts.Notifier {
	return alerts.NewMailgun(m.ApiKey, m.Domain)
}

func newInfobip(i *eventsd.Infobip) alerts.Notifier {
	return alerts.NewInfobip(i.Username,
		i.Password,
		i.ApiKey,
		i.ApplicationId,
		i.MessageId)
}

func newSlack(s *eventsd.Slack) alerts.Notifier {
	return alerts.NewSlack(s.Token, s.Channel)
}

// Watches for new vms, or vms destroyed.
func (self *User) Watch(eventsChannel *EventChannel) error {
	self.stop = make(chan struct{})
	go func() {
		for {
			select {
			case event := <-eventsChannel.channel:
				switch {
				case event.EventAction == Notify:
					err := self.alert()
					if err != nil {
						log.Warningf("Failed to process watch event: %v", err)
					}
				}
			case <-self.stop:
				log.Info("user watcher exiting")
				return
			}
		}
	}()
	return nil
}

func (self *User) Close() {
	if self.stop != nil {
		close(self.stop)
	}
}

func (self *User) alert() error {
	var err error
	for _, a := range notifiers {
		err = a.Notify("user")
	}
	return err
}
