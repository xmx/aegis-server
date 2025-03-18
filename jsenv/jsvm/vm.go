package jsvm

import "github.com/grafana/sobek"

func New() *sobek.Runtime {
	vm := sobek.New()
	mapper := sobek.TagFieldNameMapper("json", true)
	vm.SetFieldNameMapper(mapper)
	vm.SetMaxCallStackSize(64)

	return vm
}

type GlobalRegister interface {
	RegisterGlobal(vm *sobek.Runtime) error
}

func RegisterGlobals(vm *sobek.Runtime, mods []GlobalRegister) error {
	for _, mod := range mods {
		if err := mod.RegisterGlobal(vm); err != nil {
			return err
		}
	}

	return nil
}
