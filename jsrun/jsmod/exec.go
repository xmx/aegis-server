package jsmod

import (
	"os/exec"

	"github.com/xmx/aegis-server/jsrun/jsvm"
)

func NewExec() jsvm.GlobalRegister {
	return new(stdExec)
}

type stdExec struct {
	vm jsvm.Engineer
}

func (std *stdExec) RegisterGlobal(vm jsvm.Engineer) error {
	std.vm = vm
	fns := map[string]any{
		"command": std.command,
	}

	return vm.Runtime().Set("exec", fns)
}

func (std *stdExec) command(name string, args ...string) *execCommand {
	cmd := exec.Command(name, args...)
	return &execCommand{
		Cmd: cmd,
		vm:  std.vm,
	}
}

type execCommand struct {
	*exec.Cmd
	vm jsvm.Engineer
}

func (ec *execCommand) Run() error {
	ec.vm.AddFinalizer(ec.kill)
	return ec.Cmd.Run()
}

func (ec *execCommand) kill() error {
	if proc := ec.Cmd.Process; proc != nil {
		return proc.Kill()
	}
	return nil
}
