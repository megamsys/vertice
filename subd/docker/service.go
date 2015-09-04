package docker

import (
	"fmt"
	log "github.com/golang/glog"
	"github.com/megamsys/libgo/amqp"
	"github.com/megamsys/megamd/carton"
	"github.com/megamsys/megamd/meta"
	"github.com/megamsys/megamd/provision"
	"sync"
)

const QUEUE = "dockerup"

// Service manages the listener and handler for an HTTP endpoint.
type Service struct {
	wg      sync.WaitGroup
	err     chan error
	Handler *Handler

	Meta    *meta.Config
	Dockerd *Config
}

// NewService returns a new instance of Service.
func NewService(c *meta.Config, d *Config) *Service {
	s := &Service{
		err:     make(chan error),
		Meta:    c,
		Dockerd: d,
	}
	s.Handler = NewHandler(s.Dockerd)
	return s
}

// Open starts the service
func (s *Service) Open() error {
	log.Infof("Starting docker service")

	p, err := amqp.NewRabbitMQ(s.Meta.AMQP, QUEUE)
	if err != nil {
		log.Errorf("Couldn't establish an amqp (%s): %s", s.Meta, err.Error())
	}

	ch, err := p.Sub()

	for raw := range ch {
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

	carton.Provisioner = a

	if err != nil {
		//fatal(err)
		fmt.Errorf("fatal error, couldn't located the provisioner %s", s.Meta.Provider)
	}
	fmt.Printf("Using %q provisioner.\n", s.Meta.Provider)
	if initializableProvisioner, ok := carton.Provisioner.(provision.InitializableProvisioner); ok {
		err = initializableProvisioner.Initialize()
		if err != nil {
			//			fatal(err)
			fmt.Errorf("fatal error, couldn't initialize the provisioner %s", s.Meta.Provider)

		}
	}
	if messageProvisioner, ok := carton.Provisioner.(provision.MessageProvisioner); ok {
		startupMessage, err := messageProvisioner.StartupMessage()
		if err == nil && startupMessage != "" {
			fmt.Print(startupMessage)
		}
	}
}
