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
			renderer.Render(w, http.StatusUnsupportedMediaType, errorResponse{Message: "this handler can only accept json or xml"})
			return
		}

		storeRecipes, err := r.recipeLister.All()
		if err != nil {
			renderer.Render(w, http.StatusInternalServerError, errorResponse{Message: "could not fulfill your request"})
			return
		}

		var response []singleRecipeResponse
		for _, storeRecipe := range storeRecipes {
			response = append(response, singleRecipeResponse{
				ID:          strconv.FormatUint(uint64(storeRecipe.ID), 10),
				Name:        storeRecipe.Name,
				Description: storeRecipe.Description,
			})
		}

		renderer.Render(w, http.StatusOK, storeRecipes)
	})
}

// POST /recipes
func (r *Recipes) createRecipe() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Do not mix database entities with API entities
		// Bind
		// Validate
		// Logic
		r.json.Render(w, http.StatusOK, "createRecipe")
	})
}

// DELETE /recipes/{id}
func (r *Recipes) deleteRecipe() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		r.json.Render(w, http.StatusOK, "deleteRecipe")
	})
}

// UPDATE /recipes/{id}
func (r *Recipes) updateRecipe() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		r.json.Render(w, http.StatusOK, "updateRecipe")
	})
}

// GET /recipes/{id}
func (r *Recipes) readRecipe() http.HandlerFunc {
	// If the struct is only used inside one handler
	// that is fine to declare it here
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var renderer = r.json

		// Content negotiation
		// You can select the renderer and the binder
		switch r.contentNegotiator.Negotiate(req) {
		case "application/xml", "text/xml":
			renderer = r.xml
		case "application/json":
			fallthrough
		case "*/*":
			renderer = r.json
		default:
			renderer.Render(w, http.StatusUnsupportedMediaType, errorResponse{Message: "this handler can only accept json or xml"})
			return
		}

		id, err := strconv.ParseUint(r.params.ByName(req, "id"), 10, 0)
		if err != nil {
			renderer.Render(w, http.StatusBadRequest, errorResponse{Message: "id must be a positive number"})
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

			renderer.Render(w, status, errorResponse{Message: message})
			return
		}

		renderer.Render(w, http.StatusOK, singleRecipeResponse{
			ID:          strconv.FormatUint(id, 10),
			Name:        storeRecipe.Name,
			Description: storeRecipe.Description,
		})
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
			renderer.Render(w, http.StatusUnsupportedMediaType, errorResponse{Message: "this handler can only accept json or xml"})
			return
		}

		storeRecipes, err := r.recipeSearcher.Search(query)
		if err != nil {
			renderer.Render(w, http.StatusInternalServerError, errorResponse{Message: "could not fulfill your request"})
			return
		}

		var response []singleRecipeResponse
		for _, storeRecipe := range storeRecipes {
			response = append(response, singleRecipeResponse{
				ID:          strconv.FormatUint(uint64(storeRecipe.ID), 10),
				Name:        storeRecipe.Name,
				Description: storeRecipe.Description,
			})
		}

		renderer.Render(w, http.StatusOK, storeRecipes)
	})
}
