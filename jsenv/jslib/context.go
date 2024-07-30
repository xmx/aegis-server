package jslib

import (
	"context"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/xmx/aegis-server/jsenv/jsvm"
)

func Context() jsvm.Loader {
	return new(stdContext)
}

type stdContext struct {
	vm *goja.Runtime
}

func (s *stdContext) Global(vm *goja.Runtime) error {
	fields := map[string]any{
		"background":   context.Background,
		"withCancel":   s.withCancel,
		"withTimeout":  s.withTimeout,
		"withDeadline": s.withDeadline,
		"withValue":    s.withValue,
	}

	return vm.Set("context", fields)
}

func (s *stdContext) Require() (string, require.ModuleLoader) {
	return "", nil
}

func (s *stdContext) withCancel(parent context.Context) map[string]any {
	ctx, cancel := context.WithCancel(parent)
	ret := map[string]any{
		"ctx":    ctx,
		"cancel": cancel,
	}

	return ret
}

func (s *stdContext) withValue(parent context.Context, key, val any) context.Context {
	return context.WithValue(parent, key, val)
}

func (s *stdContext) withTimeout(parent context.Context, timeout time.Duration) map[string]any {
	ctx, cancel := context.WithTimeout(parent, timeout)
	ret := map[string]any{
		"ctx":    ctx,
		"cancel": cancel,
	}

	return ret
}

func (s *stdContext) withDeadline(parent context.Context, d time.Time) map[string]any {
	ctx, cancel := context.WithDeadline(parent, d)
	ret := map[string]any{
		"ctx":    ctx,
		"cancel": cancel,
	}

	return ret
}
