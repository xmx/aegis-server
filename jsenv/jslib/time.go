package jslib

import (
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/xmx/aegis-server/jsenv/jsvm"
)

func Time() jsvm.Loader {
	return new(stdTime)
}

type stdTime struct{}

func (l *stdTime) Global(vm *goja.Runtime) error {
	fields := map[string]any{
		"nanosecond":  time.Nanosecond,
		"microsecond": time.Microsecond,
		"millisecond": time.Millisecond,
		"second":      time.Second,
		"minute":      time.Minute,
		"hour":        time.Hour,
		"utc":         time.UTC,
		"local":       time.Local,
	}

	return vm.Set("time", fields)
}

func (l *stdTime) Require() (string, require.ModuleLoader) {
	return "", nil
}
