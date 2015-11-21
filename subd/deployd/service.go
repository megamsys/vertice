package deployd

import (
	"fmt"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/amqp"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megamd/carton"
	"github.com/megamsys/megamd/meta"
	"github.com/megamsys/megamd/provision"
	_ "github.com/megamsys/megamd/provision/one"
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
	c.MkGlobal() //a setter for global meta config
	return s
}

// Open starts the service
func (s *Service) Open() error {
	log.Info("starting deployd service")

	p, err := amqp.NewRabbitMQ(s.Meta.AMQP, QUEUE)
	if err != nil {
		return err
	}

	if swt, err := p.Sub(); err != nil {
		return err
	} else {
		if err = s.setProvisioner(provision.PROVIDER_ONE); err != nil {
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
	  do we have Close method in amqp ?
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
  var b map[string]string
	if initializableProvisioner, ok := carton.Provisioner.(provision.InitializableProvisioner); ok {
		err = initializableProvisioner.Initialize(s.Deployd.toMap(), b)
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
