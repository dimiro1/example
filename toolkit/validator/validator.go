package validator

type CanBeValidated interface {
	IsValid() (bool, error)
}

type Validator interface {
	Validate(v CanBeValidated) (bool, error)
}
