package jsmod

import (
	"os"

	"github.com/grafana/sobek"
	"github.com/xmx/aegis-server/jsenv/jsvm"
)

func NewOS() jsvm.GlobalRegister {
	return new(stdOS)
}

type stdOS struct {
	vm *sobek.Runtime
}

func (s *stdOS) RegisterGlobal(vm *sobek.Runtime) error {
	hm := map[string]any{
		"pid":      os.Getpid,
		"hostname": os.Hostname,
	}

	return vm.Set("os", hm)
}
