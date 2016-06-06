package events

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/events/alerts"
)

type Machine struct {
	stop chan struct{}
}

// Watches for new vms, or vms destroyed.
func (self *Machine) Watch(eventsChannel *EventChannel) error {
	self.stop = make(chan struct{})
	go func() {
		for {
			select {
			case event := <-eventsChannel.channel:
				switch {
				case event.EventAction == alerts.LAUNCHED:
					err := self.create(event)
					if err != nil {
						log.Warningf("Failed to process watch event: %v", err)
					}
				case event.EventAction == alerts.DESTROYED:
					err := self.destroy()
					if err != nil {
						log.Warningf("Failed to process watch event: %v", err)
					}
				}
			case <-self.stop:
				log.Info("machine watcher exiting")
				return
			}
		}
	}()
	return nil
}

func (self *Machine) Close() {
	if self.stop != nil {
		close(self.stop)
	}
}

func (self *Machine) create(evt *Event) error {
	return nil
}

func (self *Machine) destroy() error {
	return nil
}
