package customvalid

import (
	"context"
	"log/slog"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/xmx/aegis-server/datalayer/repository"
)

func NewValidDB(repo repository.All, log *slog.Logger) *ValidDB {
	return &ValidDB{repo: repo, log: log}
}

type ValidDB struct {
	repo repository.All
	log  *slog.Logger
}

func (vdb *ValidDB) Password() (string, validator.FuncCtx, validator.RegisterTranslationsFunc, validator.TranslationFunc) {
	const tag = "password"
	valid := func(ctx context.Context, fl validator.FieldLevel) bool {
		return false
	}
	regFunc := func(utt ut.Translator) error {
		return utt.Add(tag, "{0}不符合密码强度要求", true)
	}
	tranFun := func(ut ut.Translator, fe validator.FieldError) string {
		return "HHHHHHHHH"
	}

	return tag, valid, regFunc, tranFun
}
