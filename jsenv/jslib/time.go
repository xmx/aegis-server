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

func (l *stdTime) Global(*goja.Runtime) error {
	return nil
}

func (l *stdTime) Require() (string, require.ModuleLoader) {
	return "time", l.require
}

func (l *stdTime) require(_ *goja.Runtime, obj *goja.Object) {
	fields := map[string]any{
		"nanosecond":  time.Nanosecond,
		"microsecond": time.Microsecond,
		"millisecond": time.Millisecond,
		"second":      time.Second,
		"minute":      time.Minute,
		"hour":        time.Hour,
		"utc":         time.UTC,
		"local":       time.Local,
		"sleep":       time.Sleep,
	}
	for k, v := range fields {
		_ = obj.Set(k, v)
	}
}
