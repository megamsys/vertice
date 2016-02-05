package events

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/vertice/events/alerts"
)

const BILL = "bill"

type Config struct {
	Logo   string
	Type   string
	ApiKey string
}

type Bill struct {
	config *Config
	stop   chan struct{}
}

func NewBill(b map[string]string) *Bill {
	return &Bill{
		config: &Config{
		//	Logo:   m.Logo,
		//	Type:   e.Bill.Type,
		//	ApiKey: e.Bill.ApiKey,
		},
	}
}

// Watches for new vms, or vms destroyed.
func (self *Bill) Watch(eventsChannel *EventChannel) error {
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
				case event.EventAction == alerts.DEDUCT:
					err := self.deduct()
					if err != nil {
						log.Warningf("Failed to process watch event: %v", err)
					}
				}
			case <-self.stop:
				log.Info("bill watcher exiting")
				return
			}
		}
	}()
	return nil
}

func (self *Bill) Close() {
	if self.stop != nil {
		close(self.stop)
	}
}

func (self *Bill) create() error {
	log.Info("RECV bill create")
	return nil
}

func (self *Bill) deduct() error {
	log.Info("RECV bill deduct")
	return nil
}
