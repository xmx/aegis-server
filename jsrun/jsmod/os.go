package jsmod

import (
	"os"

	"github.com/xmx/aegis-server/jsrun/jsvm"
)

func NewOS() jsvm.ModuleRegister {
	return new(stdOS)
}

type stdOS struct{}

func (std *stdOS) RegisterModule(eng jsvm.Engineer) error {
	vals := map[string]any{
		"getpid": os.Getpid,
		"open":   os.Open,
	}
	eng.RegisterModule("os", vals, true)

	return nil
}
