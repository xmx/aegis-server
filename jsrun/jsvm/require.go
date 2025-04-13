package jsvm

import (
	"archive/zip"
	"fmt"
	"io"

	"github.com/dop251/goja"
)

type require struct {
	eng     *jsEngine
	modules map[string]goja.Value
	sources map[string]goja.Value
	source  *zip.ReadCloser
}

func (rqu *require) register(name string, module any, override bool) bool {
	_, exists := rqu.modules[name]
	if exists && !override {
		return false
	}
	rqu.modules[name] = rqu.eng.vm.ToValue(module)

	return true
}

func (rqu *require) load(call goja.FunctionCall) goja.Value {
	name := call.Argument(0).String()
	if val, exists := rqu.loadBootstrap(name); exists {
		return val
	}
	if val, exists, err := rqu.loadApplication(name); err != nil {
		panic(rqu.eng.vm.NewGoError(err))
	} else if exists {
		return val
	}

	msg := fmt.Sprintf("cannot find module '%s'", name)
	panic(rqu.eng.vm.NewTypeError(msg))
}

func (rqu *require) loadBootstrap(name string) (goja.Value, bool) {
	val, exists := rqu.modules[name]
	return val, exists
}

func (rqu *require) loadApplication(name string) (goja.Value, bool, error) {
	if rqu.source == nil {
		return nil, false, nil
	}
	if val, exists := rqu.sources[name]; exists {
		return val, true, nil
	}
	file, err := rqu.source.Open(name + ".js")
	if err != nil {
		return nil, false, err
	}
	defer file.Close()
	code, err := io.ReadAll(file)
	if err != nil {
		return nil, false, err
	}

	vm := rqu.eng.vm
	exports := vm.NewObject()
	_ = vm.Set("exports", exports)
	if _, err = rqu.eng.RunString(string(code)); err != nil {
		return nil, false, err
	}
	rqu.sources[name] = exports

	return exports, true, nil
}
