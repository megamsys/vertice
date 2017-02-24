package marketplacesd

import (
	"fmt"
	"sync"

	log "github.com/Sirupsen/logrus"
	nsq "github.com/crackcomm/nsqueue/consumer"
	"github.com/megamsys/libgo/cmd"
	constants "github.com/megamsys/libgo/utils"
	//"github.com/megamsys/vertice/carton"
	"github.com/megamsys/vertice/marketplaces"
	"github.com/megamsys/vertice/marketplaces/provision"
	_ "github.com/megamsys/vertice/marketplaces/provision/one"
	"github.com/megamsys/vertice/meta"
	"github.com/megamsys/vertice/subd/deployd"
)

const (
	TOPIC       = "marketplaces"
	maxInFlight = 150
)

// Service manages the listener and handler for an HTTP endpoint.
type Service struct {
	wg       sync.WaitGroup
	err      chan error
	Handler  *Handler
	Consumer *nsq.Consumer
	Meta     *meta.Config
	Deployd  *deployd.Config
	Config   *Config
}

// NewService returns a new instance of Service.
func NewService(c *meta.Config, config *Config, d *deployd.Config) *Service {
	s := &Service{
		err:     make(chan error),
		Meta:    c,
		Deployd: d,
		Config:  config,
	}
	s.Handler = NewHandler(s.Config)
	//c.MkGlobal() //a setter for global meta config
	return s
}

// Open starts the service
func (s *Service) Open() error {
	go func() error {
		log.Info("starting marketplacesd service")
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
	if s.Deployd.One.Enabled {
		if err := s.setProvisioner(constants.PROVIDER_ONE); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) processNSQ(msg *nsq.Message) {
	log.Debugf(TOPIC + " queue received message  :" + string(msg.Body))
	p, err := marketplaces.NewPayload(msg.Body)
	if err != nil {
		return
	}

	re, err := p.Convert()
	if err != nil {
		log.Errorf("%s", err)
		return
	}
	go s.Handler.serveNSQ(re)
	return
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

//this is an array, a property provider helps to load the provider specific stuff
func (s *Service) setProvisioner(pt string) error {
	var err error
	var tempProv provision.Provisioner

	if tempProv, err = provision.Get(pt); err != nil {
		return err
	}
	log.Debugf(cmd.Colorfy("  > configuring ", "blue", "", "bold") + fmt.Sprintf("%s ", pt))

	if initializableProvisioner, ok := tempProv.(provision.InitializableProvisioner); ok {
		err = initializableProvisioner.Initialize(s.Deployd.ToInterface(), s.Config.toMap())
		if err != nil {
			return fmt.Errorf("unable to initialize %s provisioner\n --> %s", pt, err)
		} else {
			log.Debugf(cmd.Colorfy(fmt.Sprintf("  > %s initialized", pt), "blue", "", "bold"))
		}
	}

	if messageProvisioner, ok := tempProv.(provision.MessageProvisioner); ok {
		startupMessage, err := messageProvisioner.StartupMessage()
		if err == nil && startupMessage != "" {
			log.Infof(startupMessage)
		}
	}

	marketplaces.ProvisionerMap[pt] = tempProv
	return nil
}
