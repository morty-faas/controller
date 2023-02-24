package config

import (
	"github.com/spf13/pflag"
	"net/url"
)

type Config struct {
	// Port is the listening port for the gateway server
	Port uint16
	// RIKController is the address of the RIK controller
	RIKController *url.URL
}

func NewConfig(flags *pflag.FlagSet) (Config, error) {
	var config Config

	port, err := flags.GetUint16("port")
	if err != nil {
		return config, err
	}
	config.Port = port

	rikControllerStr, err := flags.GetString("controller")
	if err != nil {
		return config, err
	}

	rikController, err := url.Parse(rikControllerStr)
	if err != nil {
		return config, err
	}
	config.RIKController = rikController

	return config, nil
}
