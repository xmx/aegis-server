package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	tzh "github.com/go-playground/validator/v10/translations/zh"
)

func TagNameFunc(tagNames []string) validator.TagNameFunc {
	size := len(tagNames)
	tags := make([]string, 0, size)
	unique := make(map[string]struct{}, len(tagNames))
	for _, tag := range tagNames {
		if tag == "" {
			continue
		}
		if _, ok := unique[tag]; ok {
			continue
		}
		unique[tag] = struct{}{}
		tags = append(tags, tag)
	}
	if len(tags) == 0 {
		return nil
	}

	return func(field reflect.StructField) string {
		var value string
		for _, tag := range tags {
			if value = field.Tag.Get(tag); value == "" || value == "-" {
				continue
			}
			if str := strings.SplitN(value, ",", 2)[0]; str != "" && str != "-" {
				value = str
				break
			}
		}

		return value
	}
}

func NewValidator(tagNameFunc validator.TagNameFunc) *Validator {
	zht := zh.New()
	ent := en.New()
	uni := ut.New(ent, zht)

	tran, _ := uni.GetTranslator(zht.Locale())
	v := validator.New()
	if tagNameFunc != nil {
		v.RegisterTagNameFunc(tagNameFunc)
	}
	_ = tzh.RegisterDefaultTranslations(v, tran)
	vd := &Validator{v: v, t: tran}

	return vd
}

type Validator struct {
	v *validator.Validate
	t ut.Translator
}

func (vd *Validator) Validate(val any) error {
	err := vd.v.Struct(val)
	switch ve := err.(type) {
	case validator.ValidationErrors:
		trans := ve.Translate(vd.t)
		return &Error{trans: trans, valid: ve}
	case *validator.InvalidValidationError:
		return &NilError{Type: ve.Type}
	}

	return err
}

func (vd *Validator) Registers(regs ...Register) error {
	for _, reg := range regs {
		if reg == nil {
			continue
		}
		tag := reg.Tag()
		if tag == "" {
			continue
		}

		if fn := reg.Validation(); fn != nil {
			if err := vd.RegisterValidation(tag, fn); err != nil {
				return err
			}
		}
		if t, tf := reg.Translation(); tf != nil {
			if err := vd.RegisterTranslation(tag, t, tf); err != nil {
				return err
			}
		}
	}

	return nil
}

func (vd *Validator) RegisterValidation(tag string, fn validator.Func) error {
	return vd.v.RegisterValidation(tag, fn)
}

func (vd *Validator) RegisterTranslation(tag string, registerFn validator.RegisterTranslationsFunc, translationFn validator.TranslationFunc) error {
	if registerFn == nil {
		return nil
	}

	return vd.v.RegisterTranslation(tag, vd.t, registerFn, translationFn)
}
