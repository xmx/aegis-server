package jsvm

import (
	"reflect"
	"strings"
	"sync"
	"unicode"

	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
)

type GlobalRegister interface {
	RegisterGlobal(vm Runtime) error
}

type Finalizer interface {
	Finalize() error
}

type Runtime interface {
	Runtime() *goja.Runtime
	RegisterGlobals(mods []GlobalRegister) error
	RunString(code string) (goja.Value, error)
	RunProgram(pgm *goja.Program) (goja.Value, error)
	Interrupt(v any)
	ClearInterrupt()
	AddFinalizer(final Finalizer)
}

func New() Runtime {
	vm := goja.New()
	vm.SetFieldNameMapper(newFieldNameMapper("json"))
	vm.SetMaxCallStackSize(64)

	return &firstRuntime{vm: vm}
}

type firstRuntime struct {
	vm    *goja.Runtime
	finmu sync.Mutex
	finhm map[Finalizer]struct{}
}

func (frt *firstRuntime) Runtime() *goja.Runtime {
	return frt.vm
}

func (frt *firstRuntime) RegisterGlobals(mods []GlobalRegister) error {
	for _, mod := range mods {
		if err := mod.RegisterGlobal(frt); err != nil {
			return err
		}
	}

	return nil
}

func (frt *firstRuntime) RunString(code string) (goja.Value, error) {
	return frt.vm.RunString(code)
}

func (frt *firstRuntime) RunProgram(pgm *goja.Program) (goja.Value, error) {
	return frt.vm.RunProgram(pgm)
}

func (frt *firstRuntime) Interrupt(v any) {
	frt.vm.Interrupt(v)

	frt.finmu.Lock()
	defer frt.finmu.Unlock()

	for final := range frt.finhm {
		_ = final.Finalize()
	}
	clear(frt.finhm)
}

func (frt *firstRuntime) ClearInterrupt() {
	frt.vm.ClearInterrupt()
}

func (frt *firstRuntime) AddFinalizer(final Finalizer) {
	if final == nil {
		return
	}

	frt.finmu.Lock()
	if frt.finhm == nil {
		frt.finhm = make(map[Finalizer]struct{}, 8)
	}
	frt.finhm[final] = struct{}{}
	frt.finmu.Unlock()
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
