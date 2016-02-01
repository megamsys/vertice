package events

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/megamd/events/alerts"
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
				case event.EventAction == alerts.CREATE:
					err := self.create()
					if err != nil {
						log.Warningf("Failed to process watch event: %v", err)
					}
				case event.EventAction == alerts.DESTROY:
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

func (self *Machine) create() error {
	log.Info("RECV machine create")
	return nil
}

func (self *Machine) destroy() error {
	log.Info("RECV machine destroy")
	return nil
}
