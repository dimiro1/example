package app

import (
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
	"github.com/dimiro1/example/toolkit/contenttype"
	"github.com/dimiro1/example/toolkit/dict"
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
	jsonBinder binder.Binder
	xmlBinder  binder.Binder

	// URL parameters extractor
	params params.ParamReader

	// Renderer
	xml  render.Renderer
	json render.Renderer

	// Detect content type
	contentType contenttype.Detector

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
	jsonBinder binder.Binder,
	xmlBinder binder.Binder,
	json render.Renderer,
	xml render.Renderer,
	contentType contenttype.Detector,
	migrator migration.Migrator,

	recipeInserter store.RecipeInserter,
	recipeFinder store.RecipeFinder,
	recipeSearcher store.RecipeSearcher,
	recipeUpdater store.RecipeUpdater,
	recipeLister store.RecipeLister) (*Application, error) {

	// make it simple to test all the parameters
	if err := anyNil(dict.M{
		"config":         config,
		"logger":         logger,
		"router":         router,
		"params":         params,
		"validator":      validator,
		"jsonBinder":     jsonBinder,
		"xmlBinder":      xmlBinder,
		"json":           json,
		"xml":            xml,
		"contentType":    contentType,
		"migrator":       migrator,
		"recipeInserter": recipeInserter,
		"recipeFinder":   recipeFinder,
		"recipeSearcher": recipeSearcher,
		"recipeUpdater":  recipeUpdater,
		"recipeLister":   recipeLister,
	}); err != nil {
		return nil, err
	}

	a := &Application{
		config: config,
		logger: logger,

		router:      router,
		params:      params,
		validator:   validator,
		jsonBinder:  jsonBinder,
		xmlBinder:   xmlBinder,
		json:        json,
		xml:         xml,
		contentType: contentType,
		migrator:    migrator,

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

// RegisterRoutes Initialize the routes
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

	if a.config.RunMigrations {
		a.logger.Info("running migrations...")
		if err := a.RunMigrations(); err != nil {
			a.logger.Error("error running migrations")
			return err
		}
		a.logger.Info("finished Running migrations...")
	}

	// Initializing routes
	a.logger.Info("registering routes...")
	a.RegisterRoutes()
	a.logger.Info("finished Registering routes...")

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

func anyNil(items dict.M) error {
	for k, v := range items {
		if v == nil {
			return fmt.Errorf("app: %s cannot be nil", k)
		}
	}
	return nil
}
