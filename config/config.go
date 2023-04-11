package config

import (
	"net/url"

	log "github.com/sirupsen/logrus"
	"github.com/thomasgouveia/go-config"
)

type Config struct {
	// Port is the listening port for the Morty controller
	Port int `yaml:"port"`
	// Cluster is the address of the RIK Controller
	Cluster string `yaml:"cluster"`
}

var loaderOptions = &config.Options[Config]{
	Format: config.YAML,

	// Environment variables lookup
	EnvEnabled: true,
	EnvPrefix:  "MORTY_CONTROLLER",

	// Configuration file
	FileName:      "controller",
	FileLocations: []string{"/etc/morty", "$HOME/.morty", "."},

	// Default configuration
	Default: &Config{
		Port:    8080,
		Cluster: "http://localhost:5000",
	},
}

// Load the configuration from the different sources (environment, files, default)
func Load() (*Config, error) {
	cl, err := config.NewLoader(loaderOptions)
	if err != nil {
		return nil, err
	}

	cfg, err := cl.Load()
	if err != nil {
		return nil, err
	}

	return cfg.validate()
}

// validate handles the configuration validation
func (c *Config) validate() (*Config, error) {
	log.Debugf("Loaded configuration: %+v", c)
	if _, err := url.Parse(c.Cluster); err != nil {
		return nil, err
	}
	return c, nil
}
