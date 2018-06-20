package app

import (
	"encoding/xml"
	"net/http"
	"strconv"

	"github.com/dimiro1/example/store"
	ct "github.com/dimiro1/example/toolkit/contenttype"
	"github.com/dimiro1/example/toolkit/render"
)

// GET /recipes
func (a *Application) listRecipes() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recipes, err := a.recipeLister.All()
		if err != nil {
			a.jsonRenderer.Render(w, http.StatusInternalServerError, err)
			return
		}
		a.jsonRenderer.Render(w, http.StatusOK, recipes)
	})
}

// POST /recipes
func (a *Application) createRecipe() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do not mix database entities with API entities
		// Bind
		// Validate
		// Logic
		a.jsonRenderer.Render(w, http.StatusOK, "createRecipe")
	})
}

// DELETE /recipes/{id}
func (a *Application) deleteRecipe() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.jsonRenderer.Render(w, http.StatusOK, "deleteRecipe")
	})
}

// UPDATE /recipes/{id}
func (a *Application) updateRecipe() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.jsonRenderer.Render(w, http.StatusOK, "updateRecipe")
	})
}

// GET /recipes/{id}
func (a *Application) readRecipe() http.HandlerFunc {
	// If the struct is only used inside one handler
	// that is fine to declare it here
	// Note: Do not use store/db entities as input or output
	// TODO: Move to it's own file
	type recipeResponse struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	// TODO: Move to it's own file
	type errorResponse struct {
		XMLName xml.Name `json:"-" xml:"error"`
		Message string   `json:"message" xml:"message,attr"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var renderer render.Renderer

		// Content negotiation
		a.contentType.Detect(r, func(t string) {
			switch t {
			case ct.XML:
				renderer = a.xmlRenderer
			case ct.JSON:
				renderer = a.jsonRenderer
			}
		})

		id, err := strconv.ParseUint(a.params.ByName(r, "id"), 10, 0)
		if err != nil {
			renderer.Render(w, http.StatusBadRequest, errorResponse{Message: "id must be a positive number"})
			return
		}

		storeRecipe, err := a.recipeFinder.Find(uint(id))
		if err != nil {
			var message string
			var status int
			if err == store.ErrRecipeNotFound {
				message = "Not Found"
				status = http.StatusNotFound
			} else {
				message = "Internal Server Error"
				status = http.StatusInternalServerError
			}

			renderer.Render(w, status, errorResponse{Message: message})
			return
		}

		renderer.Render(w, http.StatusOK, recipeResponse{
			ID:          strconv.FormatUint(id, 10),
			Name:        storeRecipe.Name,
			Description: storeRecipe.Description,
		}, nil)
	})
}

// GET /recipes/search
func (a *Application) searchRecipes() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.jsonRenderer.Render(w, http.StatusOK, "searchRecipes")
	})
}

// GET /recipes/{id}/recommendations
func (a *Application) listRecommendations() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.jsonRenderer.Render(w, http.StatusOK, "listRecommendations")
	})
}
