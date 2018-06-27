package params

import "net/http"

// ParamReader ...
type ParamReader interface {
	HasParamReader
	StringParamReader
	IntParamReader
	FloatParamReader
	BoolParamReader
}

type HasParamReader interface {
	IsPresent(r *http.Request, param string) bool
}

type StringParamReader interface {
	StringParam(r *http.Request, param string, defaultValue string) string
}

type IntParamReader interface {
	Int(r *http.Request, param string, defaultValue int) int
	Int64(r *http.Request, param string, defaultValue int64) int64

	Uint(r *http.Request, param string, defaultValue uint) uint
	Uint64(r *http.Request, param string, defaultValue uint64) uint64
}

type FloatParamReader interface {
	Float32(r *http.Request, param string, defaultValue float32) float32
	Float64(r *http.Request, param string, defaultValue float64) float64
}

type BoolParamReader interface {
	Bool(r *http.Request, param string, defaultValue bool) bool
}
