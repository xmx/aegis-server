package validation

import (
	"reflect"
	"regexp"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type Register interface {
	Tag() string
	Validation() validator.Func
	Translation() (validator.RegisterTranslationsFunc, validator.TranslationFunc)
}

func NewSemverRegister() Register {
	return &semverRegister{tag: "semver"}
}

type semverRegister struct {
	tag string
}

func (s *semverRegister) Tag() string {
	return s.tag
}

func (s *semverRegister) Validation() validator.Func {
	return nil
}

func (s *semverRegister) Translation() (validator.RegisterTranslationsFunc, validator.TranslationFunc) {
	rf := func(t ut.Translator) error {
		return t.Add(s.tag, "{0}必须是语义化版本号", true)
	}
	tf := func(t ut.Translator, fe validator.FieldError) string {
		msg, _ := t.T(s.tag, fe.Field())
		return msg
	}

	return rf, tf
}

func NewRegexRegister(tag, msg, expr string) (Register, error) {
	regex, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	reg := NewRegexpRegister(tag, msg, regex)

	return reg, nil
}

func NewRegexpRegister(tag, msg string, regex *regexp.Regexp) Register {
	fn := func(fl validator.FieldLevel) bool {
		field := fl.Field()
		if field.Kind() != reflect.String {
			return false
		}

		return regex.MatchString(field.String())
	}

	return NewFuncRegister(tag, msg, fn)
}

func NewFuncRegister(tag, msg string, fn validator.Func) Register {
	return &funcRegister{
		tag: tag,
		msg: msg,
		fn:  fn,
	}
}

type funcRegister struct {
	tag string
	msg string
	fn  validator.Func
}

func (f *funcRegister) Tag() string {
	return f.tag
}

func (f *funcRegister) Validation() validator.Func {
	return func(fl validator.FieldLevel) bool {
		return f.fn(fl)
	}
}

func (f *funcRegister) Translation() (validator.RegisterTranslationsFunc, validator.TranslationFunc) {
	rf := func(t ut.Translator) error {
		return t.Add(f.tag, f.msg, true)
	}
	tf := func(t ut.Translator, fe validator.FieldError) string {
		msg, _ := t.T(f.tag, fe.Field())
		return msg
	}

	return rf, tf
}
