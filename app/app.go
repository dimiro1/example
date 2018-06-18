package app

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/dimiro1/toolkit/router"
	"github.com/dimiro1/toolkit/render"
	"github.com/dimiro1/toolkit/validator"
	"github.com/dimiro1/toolkit/binder"
	"github.com/dimiro1/toolkit/params"

	"github.com/dimiro1/example/config"
	log "github.com/sirupsen/logrus"
)

// Application holds the application dependencies
type Application struct {
	// Application external configuration
	config *config.Config
	logger *log.Entry

	// Database interfaces/repositories

	// html templates
	// Cache interface
	// Queue interface
	// Not global tracer, like NewRelic
	// Mailer interface

	// Validates a struct
	validator validator.Validator

	// Bind struct with data from the request
	binder binder.Binder

	// URL parameters extractor
	params params.ParamReader

	// Renderer
	renderer render.Renderer

	// Mux, you are free to use any other router library
	router router.Router

	// See: https://golang.org/pkg/sync/#Once
	onceRunMigrations  sync.Once
	onceRegisterRoutes sync.Once
}

// NewApplication returns a pointer to an Application struct
func NewApplication(
	config *config.Config,
	router router.Router,
	params params.ParamReader,
	validator validator.Validator,
	binder binder.Binder,
	renderer render.Renderer,
	logger *log.Entry) (*Application, error) {

	if config == nil {
		return nil, errors.New("app: config is nil")
	}

	if router == nil {
		return nil, errors.New("app: router is nil")
	}

	if params == nil {
		return nil, errors.New("app: params is nil")
	}

	if validator == nil {
		return nil, errors.New("app: validator is nil")
	}

	if binder == nil {
		return nil, errors.New("app: binder is nil")
	}

	if renderer == nil {
		return nil, errors.New("app: renderer is nil")
	}

	if logger == nil {
		return nil, errors.New("app: logger is nil")
	}

	a := &Application{
		config:    config,
		router:    router,
		params:    params,
		validator: validator,
		binder:    binder,
		renderer:  renderer,
		logger:    logger,
	}

	return a, nil
}

// RunMigrations run needed migrations
func (a *Application) RunMigrations() error {
	a.onceRunMigrations.Do(func() {
		// TODO: Run migrations
		// Consider using github.com/dimiro1/darwin
		// you can load migrations from filesystem
		// or having directly on code
	})
	return nil
}

// Initialize the routes
func (a *Application) RegisterRoutes() error {
	// Make sure that we configure the handlers only once
	// See: https://golang.org/pkg/sync/#Once
	a.onceRegisterRoutes.Do(func() {
		// handlers returning functions, making easier to pass extra parameters
		a.router.Get("/", a.index())
	})

	return nil
}

// Start is responsible to bind to start the application
// TLS configuration, timeouts,
// It configure the routes if it was not already initialized
func (a *Application) Start() error {
	address := fmt.Sprintf(":%d", a.config.Port)

	a.logger.Infof("starting application...")

	// Running migrations if necessary
	a.logger.Info("Running migrations...")
	a.RunMigrations()

	// Initializing routes
	a.logger.Info("Registering routes...")
	a.RegisterRoutes()

	// This is the only way to safely start a http server
	// The ListenAndServe does not set timeouts
	server := &http.Server{
		Addr:         address,
		Handler:      a.router,
		ReadTimeout:  a.config.Timeouts.ReadTimeout,
		WriteTimeout: a.config.Timeouts.WriteTimeout,
		IdleTimeout:  a.config.Timeouts.IdleTimeout,
	}

	err := server.ListenAndServe()
	if err != nil {
		a.logger.WithField("address", address).Error("error serving HTTP")
		return err
	}
	return nil
}
