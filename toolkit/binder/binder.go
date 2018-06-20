package binder

import "net/http"

// Binder ...
type Binder interface {
	Bind(r *http.Request, dst interface{}) error
}
