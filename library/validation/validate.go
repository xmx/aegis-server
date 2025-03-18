package validation

import (
	"context"
	"reflect"
	"strings"

	arlocale "github.com/go-playground/locales/ar"
	enlocale "github.com/go-playground/locales/en"
	eslocale "github.com/go-playground/locales/es"
	falocale "github.com/go-playground/locales/fa"
	frlocale "github.com/go-playground/locales/fr"
	rulocale "github.com/go-playground/locales/ru"
	zhlocale "github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	artrans "github.com/go-playground/validator/v10/translations/ar"
	entrans "github.com/go-playground/validator/v10/translations/en"
	estrans "github.com/go-playground/validator/v10/translations/es"
	fatrans "github.com/go-playground/validator/v10/translations/fa"
	frtrans "github.com/go-playground/validator/v10/translations/fr"
	rutrans "github.com/go-playground/validator/v10/translations/ru"
	zhtrans "github.com/go-playground/validator/v10/translations/zh"
)

func New(tagNameFunc validator.TagNameFunc) *Validate {
	valid := validator.New()
	if tagNameFunc != nil {
		valid.RegisterTagNameFunc(tagNameFunc)
	}

	arloc := arlocale.New()
	enloc := enlocale.New()
	esloc := eslocale.New()
	faloc := falocale.New()
	frloc := frlocale.New()
	ruloc := rulocale.New()
	zhloc := zhlocale.New()

	unitran := ut.New(zhloc, arloc, enloc, esloc, faloc, frloc, ruloc, zhloc)
	artran, _ := unitran.GetTranslator(arloc.Locale())
	entran, _ := unitran.GetTranslator(enloc.Locale())
	estran, _ := unitran.GetTranslator(esloc.Locale())
	fatran, _ := unitran.GetTranslator(faloc.Locale())
	frtran, _ := unitran.GetTranslator(frloc.Locale())
	rutran, _ := unitran.GetTranslator(ruloc.Locale())
	zhtran, _ := unitran.GetTranslator(zhloc.Locale())
	trans := []ut.Translator{
		artran, entran, estran, fatran, frtran, rutran, zhtran,
	}

	_ = zhtrans.RegisterDefaultTranslations(valid, zhtran)
	_ = artrans.RegisterDefaultTranslations(valid, artran)
	_ = entrans.RegisterDefaultTranslations(valid, entran)
	_ = estrans.RegisterDefaultTranslations(valid, estran)
	_ = fatrans.RegisterDefaultTranslations(valid, fatran)
	_ = frtrans.RegisterDefaultTranslations(valid, frtran)
	_ = rutrans.RegisterDefaultTranslations(valid, rutran)

	return &Validate{
		valid:   valid,
		trans:   trans,
		unitran: unitran,
	}
}

type Validate struct {
	valid   *validator.Validate
	trans   []ut.Translator
	unitran *ut.UniversalTranslator
}

func (v *Validate) Validate(ctx context.Context, val any) error {
	err := v.valid.StructCtx(ctx, val)
	switch ve := err.(type) {
	case validator.ValidationErrors:
		return &ValidError{unitran: v.unitran, valid: ve}
	default:
		return err
	}
}

func (v *Validate) RegisterValidationCtx(tag string, fn validator.FuncCtx, callValidationEvenIfNull ...bool) error {
	return v.valid.RegisterValidationCtx(tag, fn, callValidationEvenIfNull...)
}

func (v *Validate) RegisterValidationTranslation(tag string, trans ut.Translator, registerFn validator.RegisterTranslationsFunc, translationFn validator.TranslationFunc) error {
	return v.valid.RegisterTranslation(tag, trans, registerFn, translationFn)
}

type CustomValidator interface {
	// Tag 校验器标签
	Tag() string

	// ValidationFunc 参数校验器，可以为空。
	ValidationFunc() validator.FuncCtx

	// TranslationsFunc 翻译函数。
	TranslationsFunc() validator.RegisterTranslationsFunc
}

func (v *Validate) RegisterCustomValidations(customs []CustomValidator) error {
	for _, custom := range customs {
		if err := v.RegisterCustomValidation(custom); err != nil {
			return err
		}
	}

	return nil
}

func (v *Validate) RegisterCustomValidation(custom CustomValidator) error {
	if custom == nil {
		return nil
	}
	tag := custom.Tag()
	if tag == "" {
		return nil
	}

	if validationFunc := custom.ValidationFunc(); validationFunc != nil {
		if err := v.RegisterValidationCtx(tag, validationFunc); err != nil {
			return err
		}
	}

	if translationsFunc := custom.TranslationsFunc(); translationsFunc != nil {
		for _, tran := range v.trans {
			if err := v.RegisterValidationTranslation(tag, tran, translationsFunc, v.defaultTranslation); err != nil {
				return err
			}
		}
	}

	return nil
}

func (v *Validate) defaultTranslation(utt ut.Translator, fe validator.FieldError) string {
	str, _ := utt.T(fe.Tag(), fe.Field())
	return str
}

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
