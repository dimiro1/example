package app

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/dimiro1/example/config"
	"github.com/dimiro1/example/store"

	"github.com/dimiro1/example/toolkit/binder"
	"github.com/dimiro1/example/toolkit/migration"
	"github.com/dimiro1/example/toolkit/params"
	"github.com/dimiro1/example/toolkit/render"
	"github.com/dimiro1/example/toolkit/router"
	"github.com/dimiro1/example/toolkit/validator"

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

	// Database interfaces/repositories
	// Separate into smaller interfaces is a good practice, which allows you to easily write unit tests
	recipeInserter store.RecipeInserter
	recipeFinder   store.RecipeFinder
	recipeSearcher store.RecipeSearcher
	recipeUpdater  store.RecipeUpdater
	recipeLister   store.RecipeLister

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

	// Error renderer
	// This is completely optional, your default renderer can have logic to handle errors
	errorRenderer render.Renderer

	// Database migrations
	migrator migration.Migrator

	// Mux, you are free to use any other router library
	router router.Router

	// See: https://golang.org/pkg/sync/#Once
	onceRunMigrations  sync.Once
	onceRegisterRoutes sync.Once
}

// NewApplication returns a pointer to an Application struct
func NewApplication(
	config *config.Config,
	logger *log.Entry,

	router router.Router,
	params params.ParamReader,
	validator validator.Validator,
	binder binder.Binder,
	renderer render.Renderer,
	migrator migration.Migrator,

	recipeInserter store.RecipeInserter,
	recipeFinder store.RecipeFinder,
	recipeSearcher store.RecipeSearcher,
	recipeUpdater store.RecipeUpdater,
	recipeLister store.RecipeLister) (*Application, error) {

	if config == nil {
		return nil, errors.New("app: config cannot be nil")
	}

	if logger == nil {
		return nil, errors.New("app: logger cannot be nil")
	}

	if router == nil {
		return nil, errors.New("app: router cannot be nil")
	}

	if params == nil {
		return nil, errors.New("app: params cannot be nil")
	}

	if validator == nil {
		return nil, errors.New("app: validator cannot be nil")
	}

	if binder == nil {
		return nil, errors.New("app: binder cannot be nil")
	}

	if renderer == nil {
		return nil, errors.New("app: renderer cannot be nil")
	}

	if migrator == nil {
		return nil, errors.New("app: migrator cannot be nil")
	}

	if recipeInserter == nil {
		return nil, errors.New("app: recipeInserter cannot be nil")
	}

	if recipeFinder == nil {
		return nil, errors.New("app: recipeFinder cannot be nil")
	}

	if recipeSearcher == nil {
		return nil, errors.New("app: recipeSearcher cannot be nil")
	}

	if recipeUpdater == nil {
		return nil, errors.New("app: recipeUpdater cannot be nil")
	}

	if recipeLister == nil {
		return nil, errors.New("app: recipeLister cannot be nil")
	}

	a := &Application{
		config: config,
		logger: logger,

		router:        router,
		params:        params,
		validator:     validator,
		binder:        binder,
		renderer:      renderer,
		migrator:      migrator,
		errorRenderer: renderer, // using the same renderer

		recipeInserter: recipeInserter,
		recipeFinder:   recipeFinder,
		recipeSearcher: recipeSearcher,
		recipeUpdater:  recipeUpdater,
		recipeLister:   recipeLister,
	}

	return a, nil
}

// RunMigrations run needed migrations
func (a *Application) RunMigrations() error {
	var err error
	a.onceRunMigrations.Do(func() {
		err = a.migrator.Migrate()
	})
	return err
}

// Initialize the routes
func (a *Application) RegisterRoutes() router.Router {
	// Make sure that we configure the handlers only once
	// See: https://golang.org/pkg/sync/#Once
	a.onceRegisterRoutes.Do(func() {
		// handlers returning functions, making easier to pass extra parameters
		a.router.HandleFunc("GET", "/", a.index())

		// Recipes resource
		a.router.HandleFunc("GET", "/recipes", a.listRecipes())
		a.router.HandleFunc("POST", "/recipes", a.createRecipe())
		a.router.HandleFunc("GET", "/recipes/{id:[0-9]+}", a.readRecipe())
		a.router.HandleFunc("PUT", "/recipes/{id:[0-9]+}", a.updateRecipe())
		a.router.HandleFunc("DELETE", "/recipes/{id:[0-9]+}", a.deleteRecipe())
		a.router.HandleFunc("GET", "/recipes/{id:[0-9]+}/recommendations", a.listRecommendations())
		a.router.HandleFunc("GET", "/recipes/search", a.searchRecipes())
	})

	return a.router
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
	a.logger.Info("Finished Running migrations...")

	// Initializing routes
	a.logger.Info("Registering routes...")
	a.RegisterRoutes()
	a.logger.Info("Finished Registering routes...")

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
