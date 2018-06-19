package validator

import (
	"github.com/dimiro1/example/toolkit/validator"
)

type Default struct{}

func (Default) Validate(v validator.CanBeValidated) (bool, error) {
	return v.IsValid()
}

func New() Default {
	return Default{}
}
