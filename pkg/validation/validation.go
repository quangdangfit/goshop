package validation

import (
	"fmt"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type validation struct {
	validator *validator.Validate
	uni       *ut.UniversalTranslator
	trans     *ut.Translator
}

// New Validation object
func New(opts ...Option) Validation {
	opt := getOption(opts...)
	return &validation{validator: opt.validator, uni: opt.uni, trans: opt.trans}
}

func (v *validation) ValidateStruct(s interface{}) error {
	err := v.validator.Struct(s)
	if err != nil {
		return v.Translate(err)
	}

	return nil
}

func (v *validation) Translate(err error) error {
	for _, e := range err.(validator.ValidationErrors) {
		return fmt.Errorf(e.Translate(*v.trans))
	}
	return err
}
