package jslib

import (
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/xmx/aegis-server/jsenv/jsvm"
)

func ArgsPrototype(arg any) jsvm.Loader {
	return &libArgs{
		arg: arg,
	}
}

type libArgs struct {
	arg any
}

func (l *libArgs) Global(*goja.Runtime) error {
	return nil
}

func (l *libArgs) Require() (string, require.ModuleLoader) {
	return "args", l.require
}

func (l *libArgs) require(vm *goja.Runtime, obj *goja.Object) {
	val := vm.ToValue(l.arg).ToObject(vm)
	_ = obj.SetPrototype(val)
}
