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
				case event.EventAction == alerts.DEDUCT:
					err := self.deduct(event)
					if err != nil {
						log.Warningf("Failed to process watch event: %v", err)
					}
				case event.EventAction == alerts.ONBOARD:
					err := self.OnboardFunc(event)
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

func (self *Bill) OnboardFunc(evt *Event) error {
	log.Info("onboard func")
	for _, bp := range providers {
		err := bp.Onboard(&BillOpts{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *Bill) deduct(evt *Event) error {
	log.Info("deduct bill")
	for _, bp := range providers {
		err := bp.Deduct(&BillOpts{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *Bill) Close() {
	if self.stop != nil {
		close(self.stop)
	}
}
