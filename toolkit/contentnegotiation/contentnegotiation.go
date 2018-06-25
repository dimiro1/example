package contentnegotiation

import (
	"net/http"
)

// Negotiator ...
type Negotiator interface {
	Negotiate(r *http.Request) string
}
