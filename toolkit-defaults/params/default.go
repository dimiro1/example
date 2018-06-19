package params

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Gorilla struct{}

func (Gorilla) ByName(r *http.Request, name string) string {
	vars := mux.Vars(r)
	if len(vars) == 0 {
		return ""
	}

	return vars[name]
}


func NewGorilla() Gorilla {
	return Gorilla{}
}