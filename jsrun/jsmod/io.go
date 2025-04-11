package jsmod

import (
	"io"

	"github.com/xmx/aegis-server/jsrun/jsvm"
)

func NewIO() jsvm.ModuleRegister {
	return new(stdIO)
}

type stdIO struct{}

func (std *stdIO) RegisterModule(eng jsvm.Engineer) error {
	vals := map[string]any{
		"copy":    io.Copy,
		"copyN":   io.CopyN,
		"discard": io.Discard,
	}
	eng.RegisterModule("io", vals, true)

	return nil
}
