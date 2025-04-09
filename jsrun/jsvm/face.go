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
	RegisterGlobals(mods []GlobalRegister) error
}

type finalizer interface {
	finalize() error
}
