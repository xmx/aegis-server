package jsmod

import (
	"io"

	"github.com/xmx/aegis-server/jsenv/jsvm"
)

func NewIO() jsvm.GlobalRegister {
	return new(stdIO)
}

type stdIO struct{}

func (std *stdIO) RegisterGlobal(vm jsvm.Runtime) error {
	hm := map[string]any{
		"copy":  io.Copy,
		"copyN": io.CopyN,
	}

	return vm.Runtime().Set("io", hm)
}
