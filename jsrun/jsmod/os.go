package jsmod

import (
	"io"
	"os"

	"github.com/xmx/aegis-server/jsrun/jsvm"
)

func NewOS() jsvm.GlobalRegister {
	return new(stdOS)
}

type stdOS struct{}

func (std *stdOS) RegisterGlobal(vm jsvm.Engineer) error {
	hm := map[string]any{
		"pid":    os.Getpid,
		"open":   os.Open,
		"stdout": io.Discard,
		"stderr": io.Discard,
	}

	return vm.Runtime().Set("os", hm)
}
