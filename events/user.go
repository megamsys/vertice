package events

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/megamd/subd/eventsd"
)

const (
	MAILGUN = "mailgun"
	SLACK   = "slack"
	INFOBIP = "infobip"
)

var notifiers map[string]interface{}

type User struct {
	stop chan struct{}
}

func register(e *eventsd.Config) {
	notifiers = make(map[string]interface{})
	notifiers[MAILGUN] = e.Mailgun
	notifiers[SLACK] = e.Slack
	notifiers[INFOBIP] = e.Infobip
}

func NewUser(e *eventsd.Config) *User {
	register(e)
	return &User{}
}

// Watches for new vms, or vms destroyed.
func (self *User) Watch(eventsChannel *EventChannel) error {
	self.stop = make(chan struct{})
	go func() {
		for {
			select {
			case event := <-eventsChannel.channel:
				switch {
				case event.EventAction == Alert:
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
	log.Info("RECV user alert")
	return nil
}
