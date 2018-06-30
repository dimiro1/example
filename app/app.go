package app

import (
	"fmt"
	"net/http"

	"github.com/dimiro1/example/config"
	"github.com/dimiro1/example/log"
	"github.com/dimiro1/example/toolkit/migration"
	"github.com/dimiro1/example/toolkit/module"
	"github.com/dimiro1/example/toolkit/router"
	"github.com/pkg/errors"
)

// Application holds the application dependencies
//
// You can have any kind of dependency, but, prefer to use interfaces
// it will make your life easier while writing unit tests
type Application struct {
	// Application external configuration
	config *config.Config
	logger *log.Logger

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
	logger *log.Logger,
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
			a.logger.InvalidModule(m.Name())
		}
		a.logger.RegisteringModule(m.Name())
		m.RegisterRoutes(a.router)
		a.logger.RegisteredModule(m.Name())
	}
}

// Start is responsible to bind to start the application
// TLS configuration, timeouts,
// It configure the routes if it was not already initialized
func (a *Application) Start() error {
	address := fmt.Sprintf(":%d", a.config.Port)

	a.logger.StartingApplication()

	if a.config.RunMigrations {
		a.logger.RunningMigrations()
		if err := a.RunMigrations(); err != nil {
			a.logger.ErrorRunningMigrations()
			return errors.WithStack(err)
		}
		a.logger.FinishedMigrations()
	}

	// Initializing routes
	a.logger.RegisteringRoutes()
	a.RegisterRoutes()
	a.logger.FinishedRegisteringRoutes()

	for _, route := range a.router.Routes() {
		a.logger.Route(route)
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

	a.logger.ListeningHTTP()
	err := server.ListenAndServe()
	if err != nil {
		a.logger.ErrorListeningHTTP(err, address)
		return errors.WithStack(err)
	}
	return nil
}
