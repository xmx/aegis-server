package jsmod

import (
	"context"

	"github.com/xmx/aegis-server/jsenv/jsvm"
)

func NewContext() jsvm.GlobalRegister {
	return new(stdContext)
}

type stdContext struct{}

func (*stdContext) RegisterGlobal(vm jsvm.Runtime) error {
	fns := map[string]any{
		"background":   context.Background,
		"todo":         context.TODO,
		"withCancel":   context.WithCancel,
		"withTimeout":  context.WithTimeout,
		"withValue":    context.WithValue,
		"withDeadline": context.WithDeadline,
	}

	return vm.Runtime().Set("context", fns)
}
