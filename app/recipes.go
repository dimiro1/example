package app

import (
	"net/http"
	"strconv"

	"github.com/dimiro1/example/store"
)

// GET /recipes
func (a *Application) listRecipes() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recipes, err := a.recipeLister.All()
		if err != nil {
			a.errorRenderer.Render(w, r, http.StatusInternalServerError, err)
			return
		}
		a.renderer.Render(w, r, http.StatusOK, recipes)
	})
}

// POST /recipes
func (a *Application) createRecipe() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do not mix database entities with API entities
		// Bind
		// Validate
		// Logic
		a.renderer.Render(w, r, http.StatusOK, "createRecipe")
	})
}

// DELETE /recipes/{id}
func (a *Application) deleteRecipe() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.renderer.Render(w, r, http.StatusOK, "deleteRecipe")
	})
}

// UPDATE /recipes/{id}
func (a *Application) updateRecipe() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.renderer.Render(w, r, http.StatusOK, "updateRecipe")
	})
}

// GET /recipes/{id}
func (a *Application) readRecipe() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(a.params.ByName(r, "id"), 10, 0)
		if err != nil {
			a.errorRenderer.Render(w, r, http.StatusBadRequest, "id must be a positive number")
			return
		}

		recipes, err := a.recipeFinder.Find(uint(id))
		if err == store.ErrRecipeNotFound {
			// TODO: Use a response struct
			a.errorRenderer.Render(w, r, http.StatusNotFound, "Not Found")
			return
		} else {
			a.errorRenderer.Render(w, r, http.StatusInternalServerError, err)
			return
		}
		a.renderer.Render(w, r, http.StatusOK, recipes)
	})
}

// GET /recipes/search
func (a *Application) searchRecipes() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.renderer.Render(w, r, http.StatusOK, "searchRecipes")
	})
}

// GET /recipes/{id}/recommendations
func (a *Application) listRecommendations() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.renderer.Render(w, r, http.StatusOK, "listRecommendations")
	})
}
