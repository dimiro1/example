package main

import (
	"os"

	"fmt"

	"github.com/dimiro1/example/app"
	"github.com/dimiro1/example/config"
	"github.com/dimiro1/toolkit-defaults/binder"
	"github.com/dimiro1/toolkit-defaults/params"
	"github.com/dimiro1/toolkit-defaults/render"
	"github.com/dimiro1/toolkit-defaults/router"
	"github.com/dimiro1/toolkit-defaults/validator"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Loading configs
	cfg := config.NewConfig()
	err := cfg.LoadFromEnv()
	if err != nil {
		// Log is not configured
		// Lets just call the standard panic function
		panic(err)
	}

	// Initializing log
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	if cfg.Env == "development" {
		log.SetLevel(log.DebugLevel)
	}

	logger := log.WithFields(log.Fields{
		"address": fmt.Sprintf(":%d", cfg.Port),
		"env":     cfg.Env,
	})

	// Instantiating the application
	application, err := app.NewApplication(
		cfg,
		router.New(),
		params.New(),
		validator.New(),
		binder.NewContentNegotiation(),
		render.NewContentNegotiation(),
		logger,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Running the application
	err = application.Start()
	if err != nil {
		log.Fatal(err)
	}
}
