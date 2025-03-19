package jsmod

import (
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

func (s *stdOS) RegisterGlobal(vm *goja.Runtime) error {
	hm := map[string]any{
		"pid":      os.Getpid,
		"hostname": os.Hostname,
	}

	return vm.Set("os", hm)
}
