package validator

import (
	"github.com/dimiro1/example/toolkit/validator"
)

type Simple struct{}

func (Simple) Validate(v validator.CanBeValidated) (bool, error) {
	return v.IsValid()
}

func NewSimple() Simple {
	return Simple{}
}
