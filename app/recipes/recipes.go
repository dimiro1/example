package recipes

import (
	"fmt"

	"github.com/dimiro1/example/store"
	"github.com/dimiro1/example/toolkit/binder"
	"github.com/dimiro1/example/toolkit/contentnegotiation"
	"github.com/dimiro1/example/toolkit/dict"
	"github.com/dimiro1/example/toolkit/params"
	"github.com/dimiro1/example/toolkit/render"
	"github.com/dimiro1/example/toolkit/router"
	"github.com/dimiro1/example/toolkit/validator"
	log "github.com/sirupsen/logrus"
)

// Module controller
type Recipes struct {
	logger *log.Entry

	// Database interfaces/repositories
	// Separate into smaller interfaces is a good practice, which allows you to easily write unit tests
	recipeInserter store.RecipeInserter
	recipeFinder   store.RecipeFinder
	recipeSearcher store.RecipeSearcher
	recipeUpdater  store.RecipeUpdater
	recipeLister   store.RecipeLister

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

	// Negotiate content type
	contentNegotiator contentnegotiation.Negotiator
}

func (r *Recipes) Name() string {
	return "recipes"
}

func (r *Recipes) RegisterRoutes(router router.Router) {
	router.HandleFunc("GET", "/recipes", r.listRecipes())
	router.HandleFunc("POST", "/recipes", r.createRecipe())
	router.HandleFunc("GET", "/recipes/{id:[0-9]+}", r.readRecipe())
	router.HandleFunc("GET", "/recipes/search", r.searchRecipes())
}

func NewRecipes(
	logger *log.Entry,
	params params.ParamReader,
	validator validator.Validator,
	jsonBinder binder.Binder,
	xmlBinder binder.Binder,
	json render.Renderer,
	xml render.Renderer,
	contentNegotiator contentnegotiation.Negotiator,
	recipeInserter store.RecipeInserter,
	recipeFinder store.RecipeFinder,
	recipeSearcher store.RecipeSearcher,
	recipeUpdater store.RecipeUpdater,
	recipeLister store.RecipeLister) (*Recipes, error) {

	// make it simple to test all the parameters
	if err := anyNil(dict.Dict{
		"logger":            logger,
		"params":            params,
		"validator":         validator,
		"jsonBinder":        jsonBinder,
		"xmlBinder":         xmlBinder,
		"json":              json,
		"xml":               xml,
		"contentNegotiator": contentNegotiator,
		"recipeInserter":    recipeInserter,
		"recipeFinder":      recipeFinder,
		"recipeSearcher":    recipeSearcher,
		"recipeUpdater":     recipeUpdater,
		"recipeLister":      recipeLister,
	}); err != nil {
		return nil, err
	}

	return &Recipes{
		logger:            logger,
		params:            params,
		validator:         validator,
		jsonBinder:        jsonBinder,
		xmlBinder:         xmlBinder,
		json:              json,
		xml:               xml,
		contentNegotiator: contentNegotiator,

		recipeInserter: recipeInserter,
		recipeFinder:   recipeFinder,
		recipeSearcher: recipeSearcher,
		recipeUpdater:  recipeUpdater,
		recipeLister:   recipeLister,
	}, nil
}

func anyNil(items dict.Dict) error {
	for k, v := range items {
		if v == nil {
			return fmt.Errorf("recipes: %s cannot be nil", k)
		}
	}
	return nil
}
