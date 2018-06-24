package contentnegotiation

import (
	"net/http"
	"github.com/dimiro1/example/toolkit/contentnegotiation/mediatype"
)

// Negotiator ...
type Negotiator interface {
	Negotiate(r *http.Request) mediatype.MediaType
}
