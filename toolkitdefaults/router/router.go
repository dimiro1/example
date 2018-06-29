package router

import (
	"net/http"
	"reflect"
	"runtime"

	"github.com/dimiro1/example/toolkit/router"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
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

func (r *Gorilla) Routes() []router.Route {
	var routes []router.Route

	r.mux.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error {
		methods, err := route.GetMethods()
		if err != nil {
			return errors.WithStack(err)
		}
		for _, method := range methods {
			path, err := route.GetPathTemplate()
			if err != nil {
				return errors.WithStack(err)
			}

			routes = append(routes, router.Route{
				Method:      method,
				Path:        path,
				Handler:     route.GetHandler(),
				HandlerName: runtime.FuncForPC(reflect.ValueOf(route.GetHandler()).Pointer()).Name(),
			})
		}
		return nil
	})

	return routes
}
