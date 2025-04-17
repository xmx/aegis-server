package jsvm

import (
	_ "embed"
	"sync"

	"github.com/dop251/goja"
)

var (
	//go:embed babel.js
	babelJS string

	babelPool = &sync.Pool{New: func() any {
		tsf, err := newBabel()
		return &babelCompiled{
			err: err,
			tsf: tsf,
		}
	}}
)

func Transform(code string, opts map[string]any) (string, error) {
	bc := babelPool.Get().(*babelCompiled)
	return bc.transform(code, opts)
}

type babelCompiled struct {
	err error
	tsf func(string, map[string]any) (string, error)
}

func (bc *babelCompiled) transform(code string, opts map[string]any) (string, error) {
	if err := bc.err; err != nil {
		return "", err
	}

	return bc.tsf(code, opts)
}

func newBabel() (func(string, map[string]any) (string, error), error) {
	prog, err := goja.Compile("babel.js", babelJS, false)
	if err != nil {
		return nil, err
	}

	vm := goja.New()
	logFunc := func(goja.FunctionCall) goja.Value { return nil } // babel need
	_ = vm.Set("console", map[string]func(goja.FunctionCall) goja.Value{
		"log":   logFunc,
		"error": logFunc,
		"warn":  logFunc,
	})
	if _, err = vm.RunProgram(prog); err != nil {
		return nil, err
	}

	var callable goja.Callable
	bab := vm.Get("Babel")
	if err = vm.ExportTo(bab.ToObject(vm).Get("transform"), &callable); err != nil {
		return nil, err
	}

	tsf := func(code string, opts map[string]any) (string, error) {
		if value, exx := callable(bab, vm.ToValue(code), vm.ToValue(opts)); exx != nil {
			return "", exx
		} else {
			return value.ToObject(vm).Get("code").String(), nil
		}
	}

	return tsf, nil
}
