package recipes

import (
	"encoding/xml"
	"github.com/pkg/errors"
)

// Note: Do not use store/db entities as input or output
type singleRecipe struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s singleRecipe) IsValid() (bool, error) {
	if s.Name == "" {
		return false, errors.New("name cannot be blank")
	}

	if s.Description == "" {
		return false, errors.New("description cannot be blank")
	}

	return true, nil
}

type errorResponse struct {
	XMLName xml.Name `json:"-" xml:"error"`
	Message string   `json:"message" xml:"message,attr"`
}
