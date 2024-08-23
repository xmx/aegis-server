package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/dop251/goja"
	"github.com/xmx/aegis-server/jsenv/babel"
	"github.com/xmx/aegis-server/jsenv/jsvm"
)

type GojaPlayer interface {
	Exec(ctx context.Context, loads []jsvm.Loader, script string) (goja.Value, error)
}

func NewGojaPlayer(loads []jsvm.Loader, log *slog.Logger) GojaPlayer {
	vm := jsvm.New()
	gp := &gojaPlayer{vm: vm, log: log}
	if err := jsvm.Register(vm, loads); err != nil {
		gp.err = err
	}

	return gp
}

type gojaPlayer struct {
	err error
	vm  *goja.Runtime
	log *slog.Logger
}

func (gp *gojaPlayer) Exec(ctx context.Context, loads []jsvm.Loader, script string) (goja.Value, error) {
	if err := gp.err; err != nil {
		return nil, err
	}

	vm := gp.vm
	if len(loads) != 0 {
		if err := jsvm.Register(vm, loads); err != nil {
			return nil, err
		}
	}

	commonJS, err := babel.CommonJS(script, true)
	if err != nil {
		return nil, err
	}

	timeout := time.Minute
	if deadline, ok := ctx.Deadline(); ok {
		timeout = deadline.Sub(time.Now())
	}
	timer := time.AfterFunc(timeout, func() {
		vm.Interrupt(context.DeadlineExceeded)
	})
	defer timer.Stop()

	return vm.RunString(commonJS)
}
