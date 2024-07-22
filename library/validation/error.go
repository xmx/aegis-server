package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Error struct {
	trans validator.ValidationErrorsTranslations
	valid validator.ValidationErrors
}

func (e *Error) Error() string {
	causes := make([]string, 0, len(e.valid))
	for _, err := range e.valid {
		ns := err.Namespace()
		cause := e.trans[ns]
		causes = append(causes, cause)
	}

	return strings.Join(causes, ",")
}

type NilError struct {
	Type reflect.Type
}

func (e *NilError) Error() string {
	var msg string
	if t := e.Type; t != nil {
		msg = "：" + e.Type.String()
	}

	return "无效的空类型数据" + msg
}
