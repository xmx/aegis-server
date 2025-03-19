package jsmod

import (
	"io"
	"os"

	"github.com/dop251/goja"
	"github.com/xmx/aegis-server/jsenv/jsvm"
)

func NewOS() jsvm.GlobalRegister {
	return new(stdOS)
}

type stdOS struct {
	vm *goja.Runtime
}

func (std *stdOS) RegisterGlobal(vm *goja.Runtime) error {
	std.vm = vm
	hm := map[string]any{
		"pid":    os.Getpid,
		"open":   os.Open,
		"stdout": io.Discard,
		"stderr": io.Discard,
	}

	return vm.Set("os", hm)
}
