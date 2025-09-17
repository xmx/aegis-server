package validext

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/xmx/aegis-common/validation"
)

func All() []validation.CustomValidatorFunc {
	return []validation.CustomValidatorFunc{
		mongoDB,
	}
}

func mongoDB() (string, validator.FuncCtx, validator.RegisterTranslationsFunc, validator.TranslationFunc) {
	const tag = "mongodb"
	regFunc := func(utt ut.Translator) error {
		return utt.Add(tag, "{0}不符合格式要求", true)
	}

	return tag, nil, regFunc, nil
}
