package recipes

// Note: Do not use store/db entities as input or output
type singleRecipeResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
