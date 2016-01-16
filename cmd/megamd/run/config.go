package run

import (
	"errors"

	"github.com/megamsys/megamd/meta"
	"github.com/megamsys/megamd/subd/deployd"
	"github.com/megamsys/megamd/subd/dns"
	"github.com/megamsys/megamd/subd/docker"
	"github.com/megamsys/megamd/subd/httpd"
	"github.com/megamsys/megamd/subd/metricsd"
)

type Config struct {
	Meta    *meta.Config     `toml:"meta"`
	Deployd *deployd.Config  `toml:"deployd"`
	HTTPD   *httpd.Config    `toml:"http"`
	Docker  *docker.Config   `toml:"docker"`
	Bridges *docker.Bridges  `toml:"bridges"`
	Metrics *metricsd.Config `toml:"metrics"`
	DNS     *dns.Config      `toml:"dns"`
}

func (c Config) String() string {
	return ("\n" +
		c.Meta.String() +
		c.Deployd.String() + "\n" +
		c.HTTPD.String() + "\n" +
		c.Docker.String() + "\n" +
		c.Bridges.String() + "\n" +
		c.Metrics.String() + "\n" +
		c.DNS.String())

}

// NewConfig returns an instance of Config with reasonable defaults.
func NewConfig() *Config {
	c := &Config{}
	c.Meta = meta.NewConfig()
	c.Deployd = deployd.NewConfig()
	c.HTTPD = httpd.NewConfig()
	c.Docker = docker.NewConfig()
	c.Bridges = docker.NewBridgeConfig()
	c.Metrics = metricsd.NewConfig()
	c.DNS = dns.NewConfig()
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
