package jsvm

import (
	"fmt"

	"github.com/dop251/goja"
)

type require struct {
	vm      *goja.Runtime
	modules map[string]any
}

func (rqu *require) register(name string, module any, override bool) bool {
	_, exists := rqu.modules[name]
	if exists && !override {
		return false
	}
	rqu.modules[name] = module

	return true
}

func (rqu *require) load(call goja.FunctionCall) goja.Value {
	name := call.Argument(0).String()
	if val, exists := rqu.modules[name]; exists {
		return rqu.vm.ToValue(val)
	}

	msg := fmt.Sprintf("Cannot find module '%s'", name)
	panic(rqu.vm.NewTypeError(msg))
}
