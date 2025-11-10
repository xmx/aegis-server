package validext

import (
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/xmx/aegis-common/library/validation"
)

func Customs() []validation.CustomValidatorFunc {
	return []validation.CustomValidatorFunc{
		mongoDB,
	}
}

func mongoDB() (string, validator.FuncCtx, validator.RegisterTranslationsFunc) {
	const tag = "mongodb"
	regFunc := func(utt ut.Translator) error {
		return utt.Add(tag, "{0}不符合格式要求", true)
	}

	return tag, nil, regFunc
}
