package config

import (
	"net/url"

	"github.com/polyxia-org/morty-gateway/state"
	"github.com/polyxia-org/morty-gateway/state/memory"
	"github.com/polyxia-org/morty-gateway/state/redis"
	log "github.com/sirupsen/logrus"
	"github.com/thomasgouveia/go-config"
)

type (
	Config struct {
		// Port is the listening port for the Morty controller
		Port int `yaml:"port"`
		// Cluster is the address of the RIK Controller
		Cluster string `yaml:"cluster"`
		State   State  `yaml:"state"`
	}

	State struct {
		Redis redis.Config `yaml:"redis"`
	}
)

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

// StateFactory initializes a new state implementation based on the configuration.
func (c *Config) StateFactory() (state.State, error) {
	log.Debugf("Applying state factory based on configuration")
	if err := ensureKeyHasSingleSubKey(c.State); err != nil {
		return nil, err
	}

	if isDefined(c.State.Redis) {
		return redis.NewState(&c.State.Redis)
	}

	// By default, we will use a in memory state engine if no configuration
	// is provided by the user.
	return memory.NewState(), nil
}

// validate handles the configuration validation
func (c *Config) validate() (*Config, error) {
	log.Debugf("Loaded configuration: %+v", c)
	if _, err := url.Parse(c.Cluster); err != nil {
		return nil, err
	}
	return c, nil
}
