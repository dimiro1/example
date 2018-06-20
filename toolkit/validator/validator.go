package validator

// CanBeValidated ...
type CanBeValidated interface {
	IsValid() (bool, error)
}

// Validator ...
type Validator interface {
	Validate(v CanBeValidated) (bool, error)
}
