package jsmod

import (
	"io"

	"github.com/dop251/goja"
	"github.com/xmx/aegis-server/jsenv/jsvm"
)

func NewIO() jsvm.GlobalRegister {
	return new(stdIO)
}

type stdIO struct{}

func (*stdIO) RegisterGlobal(vm *goja.Runtime) error {
	fns := map[string]any{
		"copy": io.Copy,
	}

	return vm.Set("io", fns)
}
