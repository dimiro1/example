package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Gorilla struct {
	mux *mux.Router
}

func NewGorilla() *Gorilla {
	return &Gorilla{
		mux: mux.NewRouter(),
	}
}

func (r *Gorilla) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *Gorilla) Handle(method, path string, handler http.Handler) {
	r.mux.Handle(path, handler).Methods(method)
}

func (r *Gorilla) HandleFunc(method, path string, handler http.HandlerFunc) {
	r.mux.HandleFunc(path, handler).Methods(method)
}

func (r *Gorilla) NotFound(handler http.Handler) {
	r.mux.NotFoundHandler = handler
}
