package jsvm

import (
	"github.com/dop251/goja"
)

type Engineer interface {
	Runtime() *goja.Runtime
	RunString(code string) (goja.Value, error)
	RunProgram(pgm *goja.Program) (goja.Value, error)
	RegisterModule(name string, module any, override bool) bool
	AddFinalizer(finals ...func() error)
	Interrupt(v any)
	ClearInterrupt()
}

type ModuleRegister interface {
	RegisterModule(eng Engineer) error
}

func RegisterModules(eng Engineer, mods []ModuleRegister) error {
	for _, mod := range mods {
		if mod == nil {
			continue
		}
		if err := mod.RegisterModule(eng); err != nil {
			return err
		}
	}

	return nil
}
