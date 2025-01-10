package jsvm

import "github.com/grafana/sobek"

func New() *sobek.Runtime {
	vm := sobek.New()
	mapper := sobek.TagFieldNameMapper("json", true)
	vm.SetFieldNameMapper(mapper)
	vm.SetMaxCallStackSize(64)

	return vm
}

type Module interface {
	SetGlobal(vm *sobek.Runtime) error
}

func SetModule(vm *sobek.Runtime, mods []Module) error {
	for _, mod := range mods {
		if err := mod.SetGlobal(vm); err != nil {
			return err
		}
	}

	return nil
}
