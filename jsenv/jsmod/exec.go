package jsmod

import (
	"os/exec"

	"github.com/dop251/goja"
	"github.com/xmx/aegis-server/jsenv/jsvm"
)

func NewExec() jsvm.GlobalRegister {
	return new(stdExec)
}

type stdExec struct{}

func (std *stdExec) RegisterGlobal(vm *goja.Runtime) error {
	fns := map[string]any{
		"command": exec.Command,
	}

	return vm.Set("exec", fns)
}
