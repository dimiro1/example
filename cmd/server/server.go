package main

import (
	"fmt"
	"html/template"
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

	"github.com/dimiro1/example/log"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
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
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	if cfg.Env == "development" {
		logrus.SetFormatter(&logrus.TextFormatter{})
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrusLogger := logrus.WithFields(logrus.Fields{
		"address": fmt.Sprintf(":%d", cfg.Port),
		"env":     cfg.Env,
	})

	logger := log.NewLogger(logrusLogger)

	db, err := gorm.Open("sqlite3", cfg.DatabaseDSN)
	if err != nil {
		logger.ErrorOpeningDatabase(err)
	}
	defer db.Close()

	recipeStore, err := store.NewGormRecipesStore(db)
	if err != nil {
		logger.ErrorCreatingStore(err, "GormRecipeStore")
	}

	migrator, err := store.NewGormMigrator(db)
	if err != nil {
		logger.ErrorCreatingMigrator(err)
	}

	templates, err := template.ParseGlob("templates/*")
	if err != nil {
		logger.ErrorLoadingTemplates(err, "templates/*")
	}

	homeModule, err := home.NewHome(
		logger,
		render.NewHTML(templates),
	)
	if err != nil {
		logger.ErrorInstantiatingModule(err, "home")
	}

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
		logger.ErrorInstantiatingModule(err, "recipes")
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
		logger.ErrorCreateApplication(err)
	}

	// Running the application
	err = application.Start()
	if err != nil {
		logger.ErrorStartingApplication(err)
	}
}
