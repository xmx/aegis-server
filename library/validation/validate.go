package validation

import (
	"context"
	"reflect"
	"strings"

	enlocale "github.com/go-playground/locales/en"
	zhlocale "github.com/go-playground/locales/zh"
	zhhanslocale "github.com/go-playground/locales/zh_Hans"
	zhhanttwlocate "github.com/go-playground/locales/zh_Hant_TW"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entrans "github.com/go-playground/validator/v10/translations/en"
	zhtrans "github.com/go-playground/validator/v10/translations/zh"
	zhtwtrans "github.com/go-playground/validator/v10/translations/zh_tw"
)

func New() *Validate {
	valid := validator.New()
	valid.RegisterTagNameFunc(jsonTag)

	enloc := enlocale.New()
	zhloc := zhlocale.New()
	zhhansloc := zhhanslocale.New()
	zhhanttwloc := zhhanttwlocate.New()

	unitran := ut.New(zhloc, enloc, zhloc, zhhansloc, zhhanttwloc)
	entran, _ := unitran.GetTranslator(enloc.Locale())
	zhtran, _ := unitran.GetTranslator(zhloc.Locale())
	zhhanstran, _ := unitran.GetTranslator(zhhansloc.Locale())
	zhhanttwtran, _ := unitran.GetTranslator(zhhanttwloc.Locale())
	trans := []ut.Translator{
		entran, zhtran, zhhanstran, zhhanttwtran,
	}

	_ = entrans.RegisterDefaultTranslations(valid, entran)
	_ = zhtrans.RegisterDefaultTranslations(valid, zhtran)
	_ = zhtrans.RegisterDefaultTranslations(valid, zhhanstran)
	_ = zhtwtrans.RegisterDefaultTranslations(valid, zhhanttwtran)

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

func (v *Validate) Validate(val any) error {
	err := v.valid.Struct(val)
	if ve, ok := err.(validator.ValidationErrors); ok {
		return &ValidError{unitran: v.unitran, valid: ve}
	}

	return err
}

func (v *Validate) StructCtx(ctx context.Context, val any) error {
	return v.valid.StructCtx(ctx, val)
}

func (v *Validate) RegisterValidationCtx(tag string, fn validator.FuncCtx, callValidationEvenIfNull ...bool) error {
	return v.valid.RegisterValidationCtx(tag, fn, callValidationEvenIfNull...)
}

func (v *Validate) RegisterValidationTranslation(tag string, trans ut.Translator, registerFn validator.RegisterTranslationsFunc, translationFn validator.TranslationFunc) error {
	return v.valid.RegisterTranslation(tag, trans, registerFn, translationFn)
}

type CustomValidatorFunc func() (tag string, valid validator.FuncCtx, trans validator.RegisterTranslationsFunc)

func (v *Validate) RegisterCustomValidations(customs []CustomValidatorFunc) error {
	for _, custom := range customs {
		if err := v.RegisterCustomValidation(custom); err != nil {
			return err
		}
	}

	return nil
}

func (v *Validate) RegisterCustomValidation(custom CustomValidatorFunc) error {
	if custom == nil {
		return nil
	}
	tag, validationFunc, translationsFunc := custom()
	if tag == "" {
		return nil
	}

	if validationFunc != nil {
		if err := v.RegisterValidationCtx(tag, validationFunc); err != nil {
			return err
		}
	}
	if translationsFunc != nil {
		for _, tran := range v.trans {
			if err := v.RegisterValidationTranslation(tag, tran, translationsFunc, v.defaultTranslation); err != nil {
				return err
			}
		}
	}

	return nil
}

func (v *Validate) RegisterStructValidationCtx(fn validator.StructLevelFuncCtx, types ...any) {
	v.valid.RegisterStructValidationCtx(fn, types...)
}

func (v *Validate) defaultTranslation(utt ut.Translator, fe validator.FieldError) string {
	str, _ := utt.T(fe.Tag(), fe.Field())
	return str
}

func jsonTag(f reflect.StructField) string {
	str := f.Tag.Get("json")
	cut, _, _ := strings.Cut(str, ",")
	if cut != "-" && cut != "" {
		return cut
	}

	return f.Name
}
