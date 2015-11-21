package docker

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/amqp"
	"github.com/megamsys/libgo/cmd"
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
	Bridges *Bridges
}

// NewService returns a new instance of Service.
func NewService(c *meta.Config, d *Config, b *Bridges) *Service {
	s := &Service{
		err:     make(chan error),
		Meta:    c,
		Dockerd: d,
		Bridges: b,
	}
	s.Handler = NewHandler(s.Dockerd)
	return s
}

// Open starts the service
func (s *Service) Open() error {
	log.Info("starting dockerd service")

	p, err := amqp.NewRabbitMQ(s.Meta.AMQP, QUEUE)
	if err != nil {
		return err
	}

	if swt, err := p.Sub(); err != nil {
		return err
	} else {
		if err = s.setProvisioner(provision.PROVIDER_DOCKER); err != nil {
			return err
		}
		go s.processQueue(swt)
	}

	return nil
}

// processQueue continually drains the given queue  and processes the payload
// to the appropriate request process operators.
func (s *Service) processQueue(drain chan []byte) error {
	//defer s.wg.Done()
	for raw := range drain {
		p, err := carton.NewPayload(raw)
		if err != nil {
			return err
		}

		re, err := p.Convert()
		if err != nil {
			return err
		}
		go s.Handler.serveAMQP(re)
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
func (s *Service) setProvisioner(pt string) error {
	var err error

	if carton.Provisioner, err = provision.Get(pt); err != nil {
		return err
	}
	log.Debugf(cmd.Colorfy("  > configuring ", "blue", "", "bold") + fmt.Sprintf("%s ", pt))
	if initializableProvisioner, ok := carton.Provisioner.(provision.InitializableProvisioner); ok {
		err = initializableProvisioner.Initialize(s.Dockerd.toMap(), s.Bridges.ConvertToMap())
		if err != nil {
			return fmt.Errorf("unable to initialize %s provisioner\n --> %s", pt, err)
		} else {
			log.Debugf(cmd.Colorfy(fmt.Sprintf("  > %s initialized", pt), "blue", "", "bold"))
		}
	}

	if messageProvisioner, ok := carton.Provisioner.(provision.MessageProvisioner); ok {
		startupMessage, err := messageProvisioner.StartupMessage()
		if err == nil && startupMessage != "" {
			log.Infof(startupMessage)
		}
	}
	return nil
}
