package params

import "net/http"

type ParamReader interface {
	ByName(*http.Request, string) string
}
