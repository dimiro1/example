package module

import (
	"github.com/dimiro1/example/toolkit/router"
)

type Module interface {
	RegisterRoutes(router router.Router)
	Name() string
}
