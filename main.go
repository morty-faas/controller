package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/polyxia-org/morty-gateway/config"
	"github.com/polyxia-org/morty-gateway/server"
)

const (
	logEnvVarKey = "MORTY_CONTROLLER_LOG"
)

func main() {
	// Initialize the logger before doing anything else
	level := log.InfoLevel
	envLevel := os.Getenv(logEnvVarKey)
	if envLevel != "" {
		lvl, err := log.ParseLevel(envLevel)
		if err != nil {
			log.Warnf("failed to parse log level from environment, fallback to default (INFO): %v", err)
		} else {
			level = lvl
		}
	}
	log.SetLevel(level)

	// Load the configuration for the program
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	s, err := server.NewServer(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	s.Run()
}
