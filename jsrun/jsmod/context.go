package jsmod

import (
	"context"

	"github.com/xmx/aegis-server/jsrun/jsvm"
)

func NewContext() jsvm.ModuleRegister {
	return new(stdContext)
}

type stdContext struct{}

func (*stdContext) RegisterModule(eng jsvm.Engineer) error {
	vals := map[string]any{
		"background":   context.Background,
		"todo":         context.TODO,
		"withCancel":   context.WithCancel,
		"withTimeout":  context.WithTimeout,
		"withValue":    context.WithValue,
		"withDeadline": context.WithDeadline,
	}
	eng.RegisterModule("context", vals, true)

	return nil
}
