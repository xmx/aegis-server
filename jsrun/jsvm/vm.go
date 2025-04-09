package jsvm

import (
	"reflect"
	"strings"
	"sync"
	"unicode"

	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
)

func New() Engineer {
	vm := goja.New()
	vm.SetFieldNameMapper(newFieldNameMapper("json"))
	vm.SetMaxCallStackSize(64)

	return &jsEngine{vm: vm}
}

type jsEngine struct {
	vm     *goja.Runtime
	mutex  sync.Mutex
	finals []func() error
}

func (jse *jsEngine) Runtime() *goja.Runtime {
	return jse.vm
}

func (jse *jsEngine) RunString(code string) (goja.Value, error) {
	return jse.vm.RunString(code)
}

func (jse *jsEngine) RunProgram(pgm *goja.Program) (goja.Value, error) {
	return jse.vm.RunProgram(pgm)
}

func (jse *jsEngine) AddFinalizer(finals ...func() error) {
	jse.mutex.Lock()
	defer jse.mutex.Unlock()

	for _, final := range finals {
		if final != nil {
			jse.finals = append(jse.finals, final)
		}
	}
}

func (jse *jsEngine) Interrupt(v any) {
	jse.vm.Interrupt(v)

	jse.mutex.Lock()
	defer jse.mutex.Unlock()

	for _, final := range jse.finals {
		_ = final()
	}
	clear(jse.finals)
}

func (jse *jsEngine) ClearInterrupt() {
	jse.vm.ClearInterrupt()
}

func (jse *jsEngine) RegisterGlobals(mods []GlobalRegister) error {
	for _, mod := range mods {
		if mod == nil {
			continue
		}
		if err := mod.RegisterGlobal(jse); err != nil {
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

func (fnm *fieldNameMapper) MethodName(_ reflect.Type, m reflect.Method) string {
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
