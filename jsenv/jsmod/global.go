package jsmod

import (
	"github.com/grafana/sobek"
	"github.com/xmx/aegis-server/jsenv/jsvm"
)

func NewGlobal(name string, value any) jsvm.Module {
	return &globalValue{
		name:  name,
		value: value,
	}
}

type globalValue struct {
	name  string
	value any
}

func (gv *globalValue) SetGlobal(vm *sobek.Runtime) error {
	return vm.Set(gv.name, gv.value)
}
