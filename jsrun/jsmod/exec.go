package jsmod

import (
	"os/exec"

	"github.com/xmx/aegis-server/jsrun/jsvm"
)

func NewExec() jsvm.ModuleRegister {
	return new(stdExec)
}

type stdExec struct {
	eng jsvm.Engineer
}

func (std *stdExec) RegisterModule(eng jsvm.Engineer) error {
	std.eng = eng
	vals := map[string]any{
		"command": std.command,
	}
	eng.RegisterModule("exec", vals, true)

	return nil
}

func (std *stdExec) command(name string, args ...string) *execCommand {
	cmd := exec.Command(name, args...)
	return &execCommand{
		Cmd: cmd,
		eng: std.eng,
	}
}

type execCommand struct {
	*exec.Cmd
	eng jsvm.Engineer
}

func (ec *execCommand) Run() error {
	ec.eng.AddFinalizer(ec.kill)
	return ec.Cmd.Run()
}

func (ec *execCommand) kill() error {
	if proc := ec.Cmd.Process; proc != nil {
		return proc.Kill()
	}
	return nil
}
