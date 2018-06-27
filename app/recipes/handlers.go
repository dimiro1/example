package recipes

import (
	"net/http"
	"strconv"

	"github.com/dimiro1/example/store"
)

// GET /recipes
func (r *Recipes) listRecipes() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var renderer = r.json

		// This is optional
		switch r.contentNegotiator.Negotiate(req) {
		case "application/xml", "text/xml":
			renderer = r.xml
		case "application/json":
			fallthrough
		case "*/*":
			renderer = r.json
		default:
			renderer.Render(w, req, http.StatusUnsupportedMediaType, errorResponse{Message: "this handler can only accept json or xml"})
			return
		}

		storeRecipes, err := r.recipeLister.All()
		if err != nil {
			renderer.Render(w, req, http.StatusInternalServerError, errorResponse{Message: "could not fulfill your request"})
			return
		}

		response := make([]singleRecipe, len(storeRecipes))
		for i, storeRecipe := range storeRecipes {
			response[i] = singleRecipe{
				ID:          strconv.FormatUint(uint64(storeRecipe.ID), 10),
				Name:        storeRecipe.Name,
				Description: storeRecipe.Description,
			}
		}

		if err := renderer.Render(w, req, http.StatusOK, response); err != nil {
			r.logger.WithError(err).Error("listRecipes")
		}
	})
}

// POST /recipes
func (r *Recipes) createRecipe() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var renderer = r.json
		var binder = r.jsonBinder

		// This is optional
		switch r.contentNegotiator.Negotiate(req) {
		case "application/xml", "text/xml":
			renderer = r.xml
			binder = r.xmlBinder
		case "application/json":
			fallthrough
		case "*/*":
			renderer = r.json
			binder = r.jsonBinder
		default:
			renderer.Render(w, req, http.StatusUnsupportedMediaType, errorResponse{Message: "this handler can only accept json or xml"})
			return
		}

		var input singleRecipe
		if err := binder.Bind(req, &input); err != nil {
			// TODO: better error message
			renderer.Render(w, req, http.StatusBadRequest, errorResponse{Message: "invalid input"})
			return
		}

		isValid, err := r.validator.Validate(input)
		if !isValid {
			renderer.Render(w, req, http.StatusBadRequest, errorResponse{Message: err.Error()})
			return
		}

		recipe := store.Recipe{
			Name:        input.Name,
			Description: input.Description,
		}
		if err := r.recipeInserter.Insert(&recipe); err != nil {
			r.logger.WithError(err).Error("error inserting into database")
			renderer.Render(w, req, http.StatusInternalServerError, errorResponse{Message: "error inserting into database"})
			return
		}

		if err := r.json.Render(w, req, http.StatusOK, singleRecipe{
			ID:          strconv.FormatUint(uint64(recipe.ID), 10),
			Name:        recipe.Name,
			Description: recipe.Description,
		}); err != nil {
			r.logger.WithError(err).Error("createRecipe")
		}
	})
}

// GET /recipes/{id}
func (r *Recipes) readRecipe() http.HandlerFunc {
	// If the struct is only used inside one handler
	// that is fine to declare it here
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var renderer = r.json

		switch r.contentNegotiator.Negotiate(req) {
		case "application/xml", "text/xml":
			renderer = r.xml
		case "application/json":
			fallthrough
		case "*/*":
			renderer = r.json
		default:
			renderer.Render(w, req, http.StatusUnsupportedMediaType, errorResponse{Message: "this handler can only accept json or xml"})
			return
		}

		id, err := strconv.ParseUint(r.params.ByName(req, "id"), 10, 0)
		if err != nil {
			renderer.Render(w, req, http.StatusBadRequest, errorResponse{Message: "id must be a positive number"})
			return
		}

		storeRecipe, err := r.recipeFinder.Find(uint(id))
		if err != nil {
			var (
				message = "Internal Server Error"
				status  = http.StatusInternalServerError
			)

			if err == store.ErrRecipeNotFound {
				message = "Not Found"
				status = http.StatusNotFound
			}

			renderer.Render(w, req, status, errorResponse{Message: message})
			return
		}

		if err := renderer.Render(w, req, http.StatusOK, singleRecipe{
			ID:          strconv.FormatUint(id, 10),
			Name:        storeRecipe.Name,
			Description: storeRecipe.Description,
		}); err != nil {
			r.logger.WithError(err).Error("readRecipe")
		}
	})
}

// GET /recipes/search
//TODO: Pagination, add the next page link in the response header
func (r *Recipes) searchRecipes() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var renderer = r.json
		var query = req.URL.Query().Get("q") // That is fine to use the request directly

		// This is optional
		switch r.contentNegotiator.Negotiate(req) {
		case "application/xml", "text/xml":
			renderer = r.xml
		case "application/json":
			fallthrough
		case "*/*":
			renderer = r.json
		default:
			renderer.Render(w, req, http.StatusUnsupportedMediaType, errorResponse{Message: "this handler can only accept json or xml"})
			return
		}

		storeRecipes, err := r.recipeSearcher.Search(query)
		if err != nil {
			renderer.Render(w, req, http.StatusInternalServerError, errorResponse{Message: "could not fulfill your request"})
			return
		}

		response := make([]singleRecipe, len(storeRecipes))
		for i, storeRecipe := range storeRecipes {
			response[i] = singleRecipe{
				ID:          strconv.FormatUint(uint64(storeRecipe.ID), 10),
				Name:        storeRecipe.Name,
				Description: storeRecipe.Description,
			}
		}

		if err := renderer.Render(w, req, http.StatusOK, response); err != nil {
			r.logger.WithError(err).Error("searchRecipes")
		}
	})
}
