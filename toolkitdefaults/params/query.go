package params

import (
	"net/http"
	"strconv"
)

type Query struct{}

func (Query) IsPresent(r *http.Request, param string) bool {
	values := r.URL.Query()
	_, ok := values[param]
	return ok
}

func (Query) String(r *http.Request, param string, defaultValue string) string {
	value := r.URL.Query().Get(param)
	if value == "" {
		return defaultValue
	}

	return value
}

func (q Query) Int(r *http.Request, param string, defaultValue int) int {
	return int(q.Int64(r, param, int64(defaultValue)))
}

func (q Query) Int64(r *http.Request, param string, defaultValue int64) int64 {
	value := q.String(r, param, "")
	if value == "" {
		return defaultValue
	}

	val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}

	return val
}

func (q Query) Uint(r *http.Request, param string, defaultValue uint) uint {
	return uint(q.Uint64(r, param, uint64(defaultValue)))
}

func (q Query) Uint64(r *http.Request, param string, defaultValue uint64) uint64 {
	value := q.String(r, param, "")
	if value == "" {
		return defaultValue
	}

	val, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return defaultValue
	}

	return val
}

func (q Query) Float32(r *http.Request, param string, defaultValue float32) float32 {
	return float32(q.Float64(r, param, float64(defaultValue)))
}

func (q Query) Float64(r *http.Request, param string, defaultValue float64) float64 {
	value := q.String(r, param, "")
	if value == "" {
		return defaultValue
	}

	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}

	return val
}

func (q Query) Bool(r *http.Request, param string, defaultValue bool) bool {
	value := q.String(r, param, "")
	switch value {
	case "1", "true", "t":
		return true
	}
	return defaultValue
}
