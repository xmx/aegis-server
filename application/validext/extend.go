package validext

import (
	"context"
	"net/netip"
	"reflect"

	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/xmx/aegis-common/library/validation"
	"github.com/xmx/aegis-server/library/iplist"
)

func All() []validation.CustomValidatorFunc {
	return []validation.CustomValidatorFunc{
		mongoDB,
		ipRange,
	}
}

func mongoDB() (string, validator.FuncCtx, validator.RegisterTranslationsFunc) {
	const tag = "mongodb"
	regFunc := func(utt ut.Translator) error {
		return utt.Add(tag, "{0}不是合法数据", true)
	}

	return tag, nil, regFunc
}

func ipRange() (string, validator.FuncCtx, validator.RegisterTranslationsFunc) {
	const tag = "ip_range"
	vFunc := func(ctx context.Context, fl validator.FieldLevel) bool {
		field := fl.Field()
		if field.Kind() != reflect.String {
			return false
		}

		input := field.String()
		if _, err := netip.ParseAddr(input); err == nil {
			return true
		}
		if _, err := netip.ParsePrefix(input); err == nil {
			return true
		}
		_, err := iplist.Parse([]string{input})

		return err == nil
	}
	regFunc := func(utt ut.Translator) error {
		return utt.Add(tag, "{0}不是合法的 IP、CIDR 或 IPv4 地址区间", true)
	}

	return tag, vFunc, regFunc
}
