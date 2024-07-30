package jsmod

import (
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/xmx/aegis-server/jsenv/jsvm"
)

func Args(args map[string]any) jsvm.Loader {
	return &modArg{args: args}
}

type modArg struct {
	args map[string]any
}

func (m *modArg) Global(*goja.Runtime) error {
	return nil
}

func (m *modArg) Require() (string, require.ModuleLoader) {
	return "args", m.require
}

func (m *modArg) require(_ *goja.Runtime, obj *goja.Object) {
	for k, v := range m.args {
		_ = obj.Set(k, v)
	}
}
