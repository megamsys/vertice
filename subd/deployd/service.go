package deployd

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/amqp"
	"github.com/megamsys/megamd/carton"
	"github.com/megamsys/megamd/meta"
	"github.com/megamsys/megamd/provision"

	"sync"
	"time"
)

const leaderWaitTimeout = 30 * time.Second

const QUEUE = "cloudstandup"

// Service manages the listener and handler for an HTTP endpoint.
type Service struct {
	wg      sync.WaitGroup
	err     chan error
	Handler *Handler

	Meta    *meta.Config
	Deployd *Config
}

// NewService returns a new instance of Service.
func NewService(c *meta.Config, d *Config) *Service {
	s := &Service{
		err:     make(chan error),
		Meta:    c,
		Deployd: d,
	}
	s.Handler = NewHandler(s.Deployd)
	c.MC() //an accessor.
	return s
}

// Open starts the service
func (s *Service) Open() error {
	log.Debug("starting deployd service")

	p, err := amqp.NewRabbitMQ(s.Meta.AMQP, QUEUE)
	if err != nil {
		log.Errorf("Couldn't establish an amqp (%s): %s", s.Meta.AMQP, err.Error())
	}

	drain, err := p.Sub()
	if err != nil {
		return fmt.Errorf("Couldn't subscribe to amqp (%s): %s", s.Meta.AMQP, err.Error())
	}

	s.setProvisioner()

	go s.processQueue(drain)

	return nil
}

// processQueue continually drains the given queue  and processes the queue request
// to the appropriate handlers..
func (s *Service) processQueue(drain chan []byte) error {
	//defer s.wg.Done()
	for raw := range drain {
		p, err := carton.NewPayload(raw)
		if err != nil {
			return err
		}

		pc, err := p.Convert()
		if err != nil {
			return err
		}
		go s.Handler.serveAMQP(pc)
	}
	return nil
}

// Close closes the underlying subscribe channel.
func (s *Service) Close() error {
	/*save the subscribe channel and close it.
	  don't know if the amqp has Close method ?
	  	if s.chn != nil {
	  		return s.chn.Close()
	  	}
	*/
	s.wg.Wait()
	return nil
}

// Err returns a channel for fatal errors that occur on the listener.
func (s *Service) Err() <-chan error { return s.err }

//this is an array, a property provider helps to load the provider specific stuff
func (s *Service) setProvisioner() {
	a, err := provision.Get(s.Meta.Provider)

	if err != nil {
		fmt.Errorf("fatal error, couldn't located the provisioner %s", s.Meta.Provider)
	}
	carton.Provisioner = a

	log.Debugf("Using %q provisioner. %q", s.Meta.Provider, a)
	if initializableProvisioner, ok := carton.Provisioner.(provision.InitializableProvisioner); ok {
		log.Debugf("Before initialization.")
		err = initializableProvisioner.Initialize(s.Deployd.toMap())
		if err != nil {
			log.Errorf("fatal error, couldn't initialize the provisioner %s", s.Meta.Provider)
		} else {
			log.Debugf("%s Initialized", s.Meta.Provider)
		}
	}
	log.Debugf("After initialization.")

	if messageProvisioner, ok := carton.Provisioner.(provision.MessageProvisioner); ok {
		startupMessage, err := messageProvisioner.StartupMessage()
		if err == nil && startupMessage != "" {
			log.Debugf(startupMessage)
		} else {
			log.Debugf("------> " + err.Error())
		}
	}
}
