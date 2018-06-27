package params

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Gorilla struct{}

func (g Gorilla) IsPresent(r *http.Request, param string) bool {
	values := mux.Vars(r)
	_, ok := values[param]
	return ok
}

func (g Gorilla) StringParam(r *http.Request, param string, defaultValue string) string {
	values := mux.Vars(r)
	value := values[param]

	if value == "" {
		return defaultValue
	}

	return value
}

func (g Gorilla) Int(r *http.Request, param string, defaultValue int) int {
	return int(g.Int64(r, param, int64(defaultValue)))
}

func (g Gorilla) Int64(r *http.Request, param string, defaultValue int64) int64 {
	value := g.StringParam(r, param, "")
	if value == "" {
		return defaultValue
	}

	val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}

	return val
}

func (g Gorilla) Uint(r *http.Request, param string, defaultValue uint) uint {
	return uint(g.Uint64(r, param, uint64(defaultValue)))
}

func (g Gorilla) Uint64(r *http.Request, param string, defaultValue uint64) uint64 {
	value := g.StringParam(r, param, "")
	if value == "" {
		return defaultValue
	}

	val, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return defaultValue
	}

	return val
}

func (g Gorilla) Float32(r *http.Request, param string, defaultValue float32) float32 {
	return float32(g.Float64(r, param, float64(defaultValue)))
}

func (g Gorilla) Float64(r *http.Request, param string, defaultValue float64) float64 {
	value := g.StringParam(r, param, "")
	if value == "" {
		return defaultValue
	}

	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}

	return val
}

func (g Gorilla) Bool(r *http.Request, param string, defaultValue bool) bool {
	value := g.StringParam(r, param, "")
	switch value {
	case "1", "true", "t":
		return true
	}
	return false
}
