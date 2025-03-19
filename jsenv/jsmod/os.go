package jsmod

import (
	"io"
	"os"

	"github.com/xmx/aegis-server/jsenv/jsvm"
	"github.com/xmx/aegis-server/library/machine"
)

func NewOS() jsvm.GlobalRegister {
	return new(stdOS)
}

type stdOS struct{}

func (std *stdOS) RegisterGlobal(vm jsvm.Runtime) error {
	hm := map[string]any{
		"pid":    os.Getpid,
		"open":   os.Open,
		"stdout": io.Discard,
		"stderr": io.Discard,
		"id":     machine.ID,
	}

	return vm.Runtime().Set("os", hm)
}
