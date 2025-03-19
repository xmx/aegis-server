package jsvm

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
)

func New() *goja.Runtime {
	vm := goja.New()
	vm.SetFieldNameMapper(newFieldNameMapper("json"))
	vm.SetMaxCallStackSize(64)

	return vm
}

type GlobalRegister interface {
	RegisterGlobal(vm *goja.Runtime) error
}

func RegisterGlobals(vm *goja.Runtime, mods []GlobalRegister) error {
	for _, mod := range mods {
		if err := mod.RegisterGlobal(vm); err != nil {
			return err
		}
	}

	return nil
}

func newFieldNameMapper(tagName string) goja.FieldNameMapper {
	return &fieldNameMapper{tagName: tagName}
}

type fieldNameMapper struct {
	tagName string
}

func (fnm *fieldNameMapper) FieldName(t reflect.Type, f reflect.StructField) string {
	tag := f.Tag.Get(fnm.tagName)
	if idx := strings.IndexByte(tag, ','); idx != -1 {
		tag = tag[:idx]
	}
	if parser.IsIdentifier(tag) {
		return tag
	}

	return fnm.lowerCase(f.Name)
}

func (fnm *fieldNameMapper) MethodName(t reflect.Type, m reflect.Method) string {
	return fnm.lowerCase(m.Name)
}

// lowerCase 将 Go 可导出变量转为 JS 风格的变量。
//
//	HTTP -> http
//	MyHTTP -> myHTTP
//	CopyN -> copyN
//	N -> n
func (*fieldNameMapper) lowerCase(s string) string {
	runes := []rune(s)
	size := len(runes)
	for i, r := range runes {
		if unicode.IsLower(r) {
			break
		}
		next := i + 1
		if i == 0 ||
			next >= size ||
			unicode.IsUpper(runes[next]) {
			runes[i] = unicode.ToLower(r)
		}
	}

	return string(runes)
}
