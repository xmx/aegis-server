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
	rqu := &require{
		vm:      vm,
		modules: make(map[string]any, 32),
	}
	_ = vm.Set("require", rqu.load)

	return &jsEngine{
		vm:      vm,
		require: rqu,
	}
}

type jsEngine struct {
	vm      *goja.Runtime
	require *require
	mutex   sync.Mutex
	finals  []func() error
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

func (jse *jsEngine) RegisterModule(name string, module any, override bool) bool {
	return jse.require.register(name, module, override)
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
	jse.finals = nil
}

func (jse *jsEngine) ClearInterrupt() {
	jse.vm.ClearInterrupt()
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
