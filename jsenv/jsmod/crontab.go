package jsmod

import (
	"github.com/dop251/goja"
	"github.com/xmx/aegis-server/jsenv/jsvm"
	"github.com/xmx/aegis-server/library/cronv3"
)

func NewCrontab(crond *cronv3.Crontab) jsvm.GlobalRegister {
	return &extCrontab{crond: crond}
}

type extCrontab struct {
	crond *cronv3.Crontab
}

func (ext *extCrontab) RegisterGlobal(vm *goja.Runtime) error {
	fns := map[string]any{
		"addJob": ext.crond.AddJob,
	}
	return vm.Set("crontab", fns)
}
