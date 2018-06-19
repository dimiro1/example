package binder

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

type Gorilla struct{}

func (Gorilla) Bind(r *http.Request, dst interface{}) error {
	src := map[string][]string{}

	// Route parameters
	vars := mux.Vars(r)
	for k, v := range vars {
		src[k] = append(src[k], v)
	}

	// Query
	for k, v := range r.URL.Query() {
		src[k] = append(src[k], v...)
	}

	return decoder.Decode(dst, src)
}
