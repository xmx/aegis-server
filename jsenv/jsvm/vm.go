package jsvm

import "github.com/dop251/goja"

func New() *goja.Runtime {
	vm := goja.New()
	mapper := goja.TagFieldNameMapper("json", true)
	vm.SetFieldNameMapper(mapper)
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
