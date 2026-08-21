package validation

import (
	"context"
	htmltemplate "html/template"
	"reflect"
	texttemplate "text/template"

	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func All() []CustomValidatorFunc {
	return []CustomValidatorFunc{
		uniqueFunc,
		mongodbFunc,
		requiredIfFunc,
		semverFunc,
		timezoneFunc,
		httpURLFunc,
	}
}

func uniqueFunc() (string, validator.FuncCtx, validator.RegisterTranslationsFunc) {
	const tag = "unique"
	trans := func(utt ut.Translator) error {
		return utt.Add(tag, "{0}内数据不能重复", true)
	}

	return tag, nil, trans
}

func mongodbFunc() (string, validator.FuncCtx, validator.RegisterTranslationsFunc) {
	const tag = "mongodb"
	trans := func(utt ut.Translator) error {
		return utt.Add(tag, "{0}格式错误", true)
	}

	return tag, nil, trans
}

// requiredIfFunc required_if 翻译
func requiredIfFunc() (string, validator.FuncCtx, validator.RegisterTranslationsFunc) {
	const tag = "required_if"
	trans := func(utt ut.Translator) error {
		return utt.Add(tag, "{0}为{1}时{2}必须填写", true)
	}

	return tag, nil, trans
}

// semverFunc 语义化版本号：https://semver.org/lang/zh-CN/
func semverFunc() (string, validator.FuncCtx, validator.RegisterTranslationsFunc) {
	const tag = "semver"
	trans := func(utt ut.Translator) error {
		return utt.Add(tag, "{0}必须是语义化版本号", true)
	}

	return tag, nil, trans
}

// htmlTemplateFunc html/template 模板校验
func htmlTemplateFunc() (string, validator.FuncCtx, validator.RegisterTranslationsFunc) {
	const tag = "html_template"

	valid := func(ctx context.Context, fl validator.FieldLevel) bool {
		field := fl.Field()
		if field.Kind() != reflect.String {
			return false
		}

		str := field.String()
		_, err := htmltemplate.New("").Parse(str)

		return err == nil
	}
	trans := func(utt ut.Translator) error {
		return utt.Add(tag, "{0}不符合模板语法", true)
	}

	return tag, valid, trans
}

// textTemplateFunc text/template 模板校验
func textTemplateFunc() (string, validator.FuncCtx, validator.RegisterTranslationsFunc) {
	const tag = "text_template"

	valid := func(ctx context.Context, fl validator.FieldLevel) bool {
		field := fl.Field()
		if field.Kind() != reflect.String {
			return false
		}

		str := field.String()
		_, err := texttemplate.New("").Parse(str)

		return err == nil
	}
	trans := func(utt ut.Translator) error {
		return utt.Add(tag, "{0}不符合模板语法", true)
	}

	return tag, valid, trans
}

// timezoneFunc timezone 翻译
func timezoneFunc() (string, validator.FuncCtx, validator.RegisterTranslationsFunc) {
	const tag = "timezone"
	trans := func(utt ut.Translator) error {
		return utt.Add(tag, "{0}不是合法的时区格式", true)
	}
	return tag, nil, trans
}

// httpURLFunc 翻译
func httpURLFunc() (string, validator.FuncCtx, validator.RegisterTranslationsFunc) {
	const tag = "http_url"
	trans := func(utt ut.Translator) error {
		return utt.Add(tag, "{0}格式不正确，请输入有效的 HTTP(s) 地址", true)
	}
	return tag, nil, trans
}
