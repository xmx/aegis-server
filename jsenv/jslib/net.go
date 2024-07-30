package jslib

import (
	"net"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/xmx/aegis-server/jsenv/jsvm"
)

func Net() jsvm.Loader {
	return new(stdNet)
}

type stdNet struct{}

func (n *stdNet) Global(vm *goja.Runtime) error {
	fields := map[string]any{
		"lookupIP": net.LookupIP,
	}

	return vm.Set("net", fields)
}

func (n *stdNet) Require() (string, require.ModuleLoader) {
	return "", nil
}
