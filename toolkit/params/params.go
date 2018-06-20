package params

import "net/http"

// ParamReader ...
type ParamReader interface {
	ByName(*http.Request, string) string
}
