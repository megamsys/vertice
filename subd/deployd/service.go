package deployd

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/megamsys/megamd/meta"
)

const leaderWaitTimeout = 30 * time.Second

const QUEUE = "cloudstandup"

// Service manages the listener and handler for an HTTP endpoint.
type Service struct {
	wg      sync.WaitGroup
	err     chan error
	Handler *Handler

	Meta    *meta.Config
	Deployd *deployd.Config
}

// NewService returns a new instance of Service.
func NewService(c meta.Config, d deployd.Config) (*Service, error) {
	if err != nil {
		return nil, err
	}

	s := &Service{
		err:     make(chan error),
		Meta:    &c,
		Deployd: &d,
	}
	s.Handler = NewHandler(s.Deployd)
	return s, nil
}

// Open starts the service
func (s *Service) Open() error {
	log.Infof("Starting deployd service")

	p, err := amqp.NewRabbitMQ(s.Meta.AMQP, QUEUE)
	if err != nil {
		log.Errorf("Couldn't establish an amqp (%s): %s", s.Meta, err.Error())
	}

	ch, err := p.Sub()

	for raw := range ch {
		p, err := app.NewPayload(raw)
		if err != nil {
			return err
		}
		go s.Handler.serveAMQP(p.Convert())
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
	app.Provisioner, err = provision.Get(s.Meta.Provider)
	if err != nil {
		fatal(err)
	}
	fmt.Printf("Using %q provisioner.\n", s.Provider)
	if initializableProvisioner, ok := app.Provisioner.(provision.InitializableProvisioner); ok {
		err = initializableProvisioner.Initialize()
		if err != nil {
			fatal(err)
		}
	}
	if messageProvisioner, ok := app.Provisioner.(provision.MessageProvisioner); ok {
		startupMessage, err := messageProvisioner.StartupMessage()
		if err == nil && startupMessage != "" {
			fmt.Print(startupMessage)
		}
	}
}
