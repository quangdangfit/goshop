package validation

// Validation interface
type Validation interface {
	ValidateStruct(s interface{}) error
	Translate(err error) error
}
