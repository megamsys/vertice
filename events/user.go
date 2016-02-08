package events

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/vertice/events/alerts"
)

// AfterFunc represents a after alert function, that can be registered with
// NewUser function.
type AfterFunc func(evt *Event) error

type AfterFuncs []AfterFunc

var notifiers map[string]alerts.Notifier

type User struct {
	stop chan struct{}
	fns  AfterFuncs
}

func NewUser(e EventsConfigMap, fns AfterFuncs) *User {
	register(e)
	return &User{fns: fns}
}

func register(e EventsConfigMap) {
	notifiers = make(map[string]alerts.Notifier)
	notifiers[alerts.MAILGUN] = newMailgun(e.Get(alerts.MAILGUN))
	notifiers[alerts.INFOBIP] = newInfobip(e.Get(alerts.INFOBIP))
	notifiers[alerts.SLACK] = newSlack(e.Get(alerts.SLACK))
}

func newMailgun(m map[string]string) alerts.Notifier {
	return alerts.NewMailgun(m)
}

func newInfobip(m map[string]string) alerts.Notifier {
	return alerts.NewInfobip(m)
}

func newSlack(m map[string]string) alerts.Notifier {
	return alerts.NewSlack(m)
}

// Watches for new vms, or vms destroyed.
func (self *User) Watch(eventsChannel *EventChannel) error {
	self.stop = make(chan struct{})
	go func() {
		for {
			select {
			case event := <-eventsChannel.channel:
				err := self.alert(event)
				if err != nil {
					log.Warningf("Failed to process watch event: %v", err)
				}
			case <-self.stop:
				log.Info("user watcher exiting")
				return
			}
		}
	}()
	return nil
}

func (self *User) alert(evt *Event) error {
	var err error
	for _, a := range notifiers {
		err = a.Notify(evt.EventAction, evt.EventData.M)
	}
	if err != nil {
		return err
	}
	return self.after(evt)
}

func (self *User) after(evt *Event) error {
	var err error
	for _, fn := range self.fns {
		err = fn(evt)
	}
	return err
}

func (self *User) Close() {
	if self.stop != nil {
		close(self.stop)
	}
}
