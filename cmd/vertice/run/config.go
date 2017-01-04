package run

import (
	"errors"

	"github.com/megamsys/vertice/meta"
	"github.com/megamsys/vertice/storage"
	"github.com/megamsys/vertice/subd/deployd"
	"github.com/megamsys/vertice/subd/dns"
	"github.com/megamsys/vertice/subd/docker"
	"github.com/megamsys/vertice/subd/eventsd"
	"github.com/megamsys/vertice/subd/httpd"
	"github.com/megamsys/vertice/subd/metricsd"
)

type Config struct {
	Meta    *meta.Config     `toml:"meta"`
	Deployd *deployd.Config  `toml:"deployd"`
	HTTPD   *httpd.Config    `toml:"http"`
	Docker  *docker.Config   `toml:"docker"`
	Metrics *metricsd.Config `toml:"metrics"`
	DNS     *dns.Config      `toml:"dns"`
	Events  *eventsd.Config  `toml:"events"`
	Storage *storage.Config  `toml:"storage"`
}

func (c Config) String() string {
	return ("\n" +
		c.Meta.String() +
		c.Deployd.String() + "\n" +
		c.HTTPD.String() + "\n" +
		c.Docker.String() + "\n" +
		c.Metrics.String() + "\n" +
		c.DNS.String() + "\n" +
		c.Events.String() + "\n" +
	  c.Storage.String())
}

// NewConfig returns an instance of Config with reasonable defaults.
func NewConfig() *Config {
	c := &Config{}
	c.Meta = meta.NewConfig()
	c.Deployd = deployd.NewConfig()
	c.HTTPD = httpd.NewConfig()
	c.Docker = docker.NewConfig()
	c.Metrics = metricsd.NewConfig()
	c.Events = eventsd.NewConfig()
	c.DNS = dns.NewConfig()
	c.Storage = storage.NewConfig()
	return c
}

// NewDemoConfig returns the config that runs when no config is specified.
func NewDemoConfig() (*Config, error) {
	c := NewConfig()
	return c, nil
}

// Validate returns an error if the config is invalid.
func (c *Config) Validate() error {
	if c.Meta.Dir == "" {
		return errors.New("Meta.Dir must be specified")
	}
	return nil
}
