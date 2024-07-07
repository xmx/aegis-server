package jsmod

import (
	"os"
	"strings"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/xmx/aegis-server/jsenv/jsvm"
)

func OS() jsvm.Loader {
	return new(stdOS)
}

type stdOS struct{}

func (o *stdOS) Global(vm *goja.Runtime) error {
	fields := map[string]any{
		"args":     os.Args,
		"getwd":    os.Getwd,
		"getpid":   os.Getpid,
		"environ":  o.environ(),
		"hostname": os.Hostname,
	}

	return vm.Set("os", fields)
}

func (o *stdOS) Require() (string, require.ModuleLoader) {
	return "", nil
}

func (o *stdOS) environ() map[string]string {
	envs := os.Environ()
	hms := make(map[string]string, len(envs))
	for _, env := range envs {
		sn := strings.SplitN(env, "=", 2)
		if len(sn) != 2 || sn[0] == "" {
			continue
		}
		hms[sn[0]] = sn[1]
	}

	return hms
}
