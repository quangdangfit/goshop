package validation

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	enLocales "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/quangdangfit/gocommon/logger"
)

// PhoneNumberRegex constants
const (
	PhoneNumberRegex = "^[0-9]{7,14}$"
)

// Option validation option
type Option interface {
	apply(*option)
}

// option implement
type option struct {
	validator *validator.Validate
	uni       *ut.UniversalTranslator
	trans     *ut.Translator
}

type optionFn func(*option)

func (optFn optionFn) apply(opt *option) {
	optFn(opt)
}

// WithValidator set validator
func WithValidator(v *validator.Validate) Option {
	return optionFn(func(opt *option) {
		opt.validator = v
	})
}

// WithUniversalTranslator set UniversalTranslator
func WithUniversalTranslator(uni *ut.UniversalTranslator) Option {
	return optionFn(func(opt *option) {
		opt.uni = uni
	})
}

// WithTranslator set Translator
func WithTranslator(trans *ut.Translator) Option {
	return optionFn(func(opt *option) {
		opt.trans = trans
	})
}

func getDefaultOption() option {
	v := validator.New()

	translator := enLocales.New()
	uni := ut.New(translator, translator)

	trans, found := uni.GetTranslator("en")
	if !found {
		logger.Error("translator not found")
	}

	if err := enTranslations.RegisterDefaultTranslations(v, trans); err != nil {
		logger.Error(err)
	}

	_ = v.RegisterTranslation("password", trans, func(ut ut.Translator) error {
		return ut.Add("password", "{0} is not strong enough, password must be at least 6 characters", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("password", fe.Field())
		return t
	})

	_ = v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		pwdLen := len(fl.Field().String())
		if pwdLen != 0 && (pwdLen < 6) {
			return false
		}
		return true
	})

	_ = v.RegisterTranslation("countryCode", trans, func(ut ut.Translator) error {
		return ut.Add("countryCode", "{0} must be at least 2 characters and start with '+'", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("countryCode", fe.Field())
		return t
	})

	_ = v.RegisterValidation("countryCode", func(fl validator.FieldLevel) bool {
		codeLen := len(fl.Field().String())
		if codeLen != 0 && (codeLen < 2 || !strings.HasPrefix(fl.Field().String(), "+")) {
			return false
		}
		return true
	})

	_ = v.RegisterTranslation("phone", trans, func(ut ut.Translator) error {
		return ut.Add("phone", "{0} must be valid phone", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("phone", fe.Field())
		return t
	})

	_ = v.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		_, err := strconv.Atoi(fl.Field().String())
		if err != nil {
			return false
		}

		matched, err := regexp.MatchString(PhoneNumberRegex, fl.Field().String())
		if !matched || err != nil {
			return false
		}

		return true
	})

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		jsonTag := fld.Tag.Get("json")
		if jsonTag == "" {
			return fld.Name
		}

		name := strings.SplitN(jsonTag, ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})

	return option{
		validator: v,
		uni:       uni,
		trans:     &trans,
	}
}

func getOption(opts ...Option) option {
	opt := getDefaultOption()
	for _, o := range opts {
		o.apply(&opt)
	}

	return opt
}
