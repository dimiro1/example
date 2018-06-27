package main

import (
	"fmt"
	"os"

	"github.com/dimiro1/example/app"
	"github.com/dimiro1/example/app/home"
	"github.com/dimiro1/example/app/recipes"
	"github.com/dimiro1/example/config"
	"github.com/dimiro1/example/store"
	"github.com/dimiro1/example/toolkitdefaults/binder"
	"github.com/dimiro1/example/toolkitdefaults/contentnegotiation"
	"github.com/dimiro1/example/toolkitdefaults/params"
	"github.com/dimiro1/example/toolkitdefaults/render"
	"github.com/dimiro1/example/toolkitdefaults/router"
	"github.com/dimiro1/example/toolkitdefaults/validator"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Loading configs
	cfg, err := config.FromEnv()
	if err != nil {
		// Log is not configured
		// Lets just call the standard panic function
		panic(err)
	}

	// Initializing log
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.JSONFormatter{})

	if cfg.Env == "development" {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.DebugLevel)
	}

	logger := log.WithFields(log.Fields{
		"address": fmt.Sprintf(":%d", cfg.Port),
		"env":     cfg.Env,
	})

	db, err := gorm.Open("sqlite3", cfg.DatabaseDSN)
	if err != nil {
		log.WithError(err).Fatal(err)
	}
	defer db.Close()

	recipeStore, err := store.NewGormRecipesStore(db)
	if err != nil {
		log.WithError(err).Fatal("failed to create recipe store")
	}

	migrator, err := store.NewGormMigrator(db)
	if err != nil {
		log.WithError(err).Fatal("failed to create db migrator")
	}

	homeModule := home.NewHome(render.Text{})
	recipesModule, err := recipes.NewRecipes(
		logger,
		params.Gorilla{},
		params.Query{},
		validator.NewBasic(),
		binder.JSON{},
		binder.XML{},
		render.JSON{},
		render.XML{},
		contentnegotiation.NewNegotiator(
			contentnegotiation.Offers("application/json", "application/xml", "text/xml"),
		),
		// Stores
		recipeStore,
		recipeStore,
		recipeStore,
		recipeStore,
		recipeStore,
	)
	if err != nil {
		log.WithError(err).Fatal("failed to create new recipes module")
	}

	// Instantiating the application
	application, err := app.NewApplication(
		cfg,
		logger,
		router.NewGorilla(),
		migrator,

		// modules
		homeModule,
		recipesModule,
	)
	if err != nil {
		log.WithError(err).Fatal("failed to create new application")
	}

	// Running the application
	err = application.Start()
	if err != nil {
		log.WithError(err).Fatal("failed to start the application")
	}
}
