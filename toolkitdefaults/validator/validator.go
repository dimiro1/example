package validator

import (
	"github.com/dimiro1/example/toolkit/validator"
)

type Basic struct{}

func (Basic) Validate(v validator.CanBeValidated) (bool, error) {
	return v.IsValid()
}

func NewBasic() Basic {
	return Basic{}
}
