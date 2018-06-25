package app

import (
	"net/http"
	"strconv"

	"github.com/dimiro1/example/store"
	"github.com/dimiro1/example/toolkit/mediatype"
)

// GET /recipes
func (a *Application) listRecipes() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var renderer = a.json

		// This is optional
		switch a.contentNegotiator.Negotiate(r) {
		case mediatype.ApplicationXML:
			renderer = a.xml
		case mediatype.ApplicationJSON:
			fallthrough
		case mediatype.All:
			renderer = a.json
		default:
			renderer.Render(w, http.StatusUnsupportedMediaType, errorResponse{Message: "this handler can only accept json or xml"})
			return
		}

		storeRecipes, err := a.recipeLister.All()
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
func (a *Application) createRecipe() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do not mix database entities with API entities
		// Bind
		// Validate
		// Logic
		a.json.Render(w, http.StatusOK, "createRecipe")
	})
}

// DELETE /recipes/{id}
func (a *Application) deleteRecipe() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.json.Render(w, http.StatusOK, "deleteRecipe")
	})
}

// UPDATE /recipes/{id}
func (a *Application) updateRecipe() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.json.Render(w, http.StatusOK, "updateRecipe")
	})
}

// GET /recipes/{id}
func (a *Application) readRecipe() http.HandlerFunc {
	// If the struct is only used inside one handler
	// that is fine to declare it here
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var renderer = a.json

		// Content negotiation
		// You can select the renderer and the binder
		switch a.contentNegotiator.Negotiate(r) {
		case mediatype.ApplicationXML:
			renderer = a.xml
		case mediatype.ApplicationJSON:
			fallthrough
		case mediatype.All:
			renderer = a.json
		default:
			renderer.Render(w, http.StatusUnsupportedMediaType, errorResponse{Message: "this handler can only accept json or xml"})
			return
		}

		id, err := strconv.ParseUint(a.params.ByName(r, "id"), 10, 0)
		if err != nil {
			renderer.Render(w, http.StatusBadRequest, errorResponse{Message: "id must be a positive number"})
			return
		}

		storeRecipe, err := a.recipeFinder.Find(uint(id))
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
func (a *Application) searchRecipes() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var renderer = a.json
		var query = r.URL.Query().Get("q") // That is fine to use the request directly

		// This is optional
		switch a.contentNegotiator.Negotiate(r) {
		case mediatype.ApplicationXML:
			renderer = a.xml
		case mediatype.ApplicationJSON:
			fallthrough
		case mediatype.All:
			renderer = a.json
		default:
			renderer.Render(w, http.StatusUnsupportedMediaType, errorResponse{Message: "this handler can only accept json or xml"})
			return
		}

		storeRecipes, err := a.recipeSearcher.Search(query)
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
