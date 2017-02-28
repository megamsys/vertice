package eventsd

import (
	log "github.com/Sirupsen/logrus"
	nsq "github.com/crackcomm/nsqueue/consumer"
	"github.com/megamsys/libgo/events"
	"github.com/megamsys/vertice/meta"
	"github.com/megamsys/vertice/subd/deployd"
	"sync"
)

const (
	TOPIC       = "events"
	maxInFlight = 150
)

// Service manages the listener and handler for an HTTP endpoint.
type Service struct {
	wg       sync.WaitGroup
	err      chan error
	Handler  *Handler
	Consumer *nsq.Consumer
	Meta     *meta.Config
	Eventsd  *Config
}

// NewService returns a new instance of Service.
func NewService(c *meta.Config, e *Config, d *deployd.Config) *Service {
	s := &Service{
		err:     make(chan error),
		Meta:    c,
		Eventsd: e,
	}
	s.Handler = NewHandler(nil)
	s.Handler.Deployd = d
	return s
}

// Open starts the service
func (s *Service) Open() error {
	if err := s.setEventsWrap(s.Eventsd); err != nil {
		return err
	}
	go func() error {
		log.Info("starting eventsd service")
		if err := nsq.Register(TOPIC, "engine", maxInFlight, s.processNSQ); err != nil {
			return err
		}
		if err := nsq.Connect(s.Meta.NSQd...); err != nil {
			return err
		}

		s.Consumer = nsq.DefaultConsumer

		nsq.Start(true)
		return nil
	}()
	return nil
}

func (s *Service) processNSQ(msg *nsq.Message) {
	log.Debugf(TOPIC + " queue received message  :" + string(msg.Body))
	pe, err := events.NewParseEvent(msg.Body)
	if err != nil {
		return
	}

	e, err := pe.AsEvent()
	if err != nil {
		return
	}

	go s.Handler.serveNSQ(e, pe.Inputs.Match("email"))
	return
}

func (s *Service) setEventsWrap(e *Config) error {
	return events.NewWrap(e.toMap())
}

// Close closes the underlying subscribe channel.
func (s *Service) Close() error {
	if s.Consumer != nil {
		s.Consumer.Stop()
	}

	s.wg.Wait()
	return nil
}

// Err returns a channel for fatal errors that occur on the listener.
func (s *Service) Err() <-chan error { return s.err }
