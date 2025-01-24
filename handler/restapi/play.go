package restapi

import (
	"context"
	"io"
	"time"

	"github.com/xmx/aegis-server/jsenv/jsmod"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/jsenv/jsvm"
)

func NewPlay(mods []jsvm.Module) *Play {
	return &Play{
		mods: mods,
	}
}

type Play struct {
	mods []jsvm.Module
}

func (p *Play) Route(r *ship.RouteGroupBuilder) error {
	r.Route("/play/js").POST(p.run)
	return nil
}

func (p *Play) run(c *ship.Context) error {
	vm := jsvm.New()
	timer := time.AfterFunc(time.Hour, func() {
		vm.Interrupt(context.DeadlineExceeded)
	})
	defer timer.Stop()

	console := jsmod.NewConsole(c)
	console.SetGlobal(vm)
	mods := p.mods
	if err := jsvm.RegisterModules(vm, "server", mods); err != nil {
		return err
	}

	str, _ := io.ReadAll(c.Body())
	_, err := vm.RunString(string(str))
	if err != nil {
		return err
	}

	return nil
}
