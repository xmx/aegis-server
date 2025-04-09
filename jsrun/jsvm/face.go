package jsvm

import "github.com/dop251/goja"

type GlobalRegister interface {
	RegisterGlobal(eng Engineer) error
}

type Engineer interface {
	Runtime() *goja.Runtime
	RunString(code string) (goja.Value, error)
	RunProgram(pgm *goja.Program) (goja.Value, error)
	AddFinalizer(finals ...func() error)
	Interrupt(v any)
	ClearInterrupt()
}

func RegisterGlobals(eng Engineer, mods []GlobalRegister) error {
	for _, mod := range mods {
		if mod == nil {
			continue
		}
		if err := mod.RegisterGlobal(eng); err != nil {
			return err
		}
	}

	return nil
}
