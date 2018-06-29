package app

import (
	"fmt"
	"net/http"

	"github.com/dimiro1/example/config"
	"github.com/dimiro1/example/toolkit/migration"
	"github.com/dimiro1/example/toolkit/module"
	"github.com/dimiro1/example/toolkit/router"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Application holds the application dependencies
//
// You can have any kind of dependency, but, prefer to use interfaces
// it will make your life easier while writing unit tests
type Application struct {
	// Application external configuration
	config *config.Config
	logger *log.Entry

	// Database migrations
	migrator migration.Migrator

	// Mux, you are free to use any other router library
	router router.Router

	// modules modules attached to this application
	modules []module.Module
}

// NewApplication returns a pointer to an Application struct
func NewApplication(
	config *config.Config,
	logger *log.Entry,
	router router.Router,
	migrator migration.Migrator,
	modules ...module.Module) (*Application, error) {

	if config == nil {
		return nil, errors.New("app: config cannot be nil")
	}

	if logger == nil {
		return nil, errors.New("app: logger cannot be nil")
	}

	if router == nil {
		return nil, errors.New("app: router cannot be nil")
	}

	if migrator == nil {
		return nil, errors.New("app: migrator cannot be nil")
	}

	a := &Application{
		config:   config,
		logger:   logger,
		router:   router,
		migrator: migrator,
		modules:  modules,
	}

	return a, nil
}

// RunMigrations run needed migrations
func (a *Application) RunMigrations() error {
	return errors.WithStack(a.migrator.Migrate())
}

// RegisterRoutes Initialize the routes
func (a *Application) RegisterRoutes() {
	for _, m := range a.modules {
		if m == nil {
			a.logger.WithField("name", m.Name()).Error("could not register module")
		}
		a.logger.WithField("name", m.Name()).Debug("registering module...")
		m.RegisterRoutes(a.router)
		a.logger.WithField("name", m.Name()).Debug("module registered...")
	}
}

// Start is responsible to bind to start the application
// TLS configuration, timeouts,
// It configure the routes if it was not already initialized
func (a *Application) Start() error {
	address := fmt.Sprintf(":%d", a.config.Port)

	a.logger.Infof("starting application...")

	if a.config.RunMigrations {
		a.logger.Debug("running migrations...")
		if err := a.RunMigrations(); err != nil {
			a.logger.Error("error running migrations")
			return errors.WithStack(err)
		}
		a.logger.Debug("finished Running migrations...")
	}

	// Initializing routes
	a.logger.Debug("registering routes...")
	a.RegisterRoutes()
	a.logger.Debug("finished Registering routes...")

	a.logger.Info("routes registered...")
	for _, route := range a.router.Routes() {
		a.logger.WithFields(log.Fields{
			"method":  route.Method,
			"route":   route.Path,
			"handler": route.HandlerName,
		}).Info()
	}

	// This is the only way to safely start a http server
	// The ListenAndServe does not set timeouts
	server := &http.Server{
		Addr:         address,
		Handler:      a.router,
		ReadTimeout:  a.config.Timeouts.ReadTimeout,
		WriteTimeout: a.config.Timeouts.WriteTimeout,
		IdleTimeout:  a.config.Timeouts.IdleTimeout,
	}

	a.logger.Info("listening...")
	err := server.ListenAndServe()
	if err != nil {
		a.logger.WithField("address", address).Error("error serving HTTP")
		return errors.WithStack(err)
	}
	return nil
}
