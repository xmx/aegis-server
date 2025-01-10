package restapi

import (
	"context"
	"io"
	"time"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/jsenv/jsmod"
	"github.com/xmx/aegis-server/jsenv/jsvm"
	"github.com/xmx/aegis-server/profile"
)

func NewPlay(cfg *profile.Config) *Play {
	return &Play{
		cfg: cfg,
	}
}

type Play struct {
	cfg *profile.Config
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

	mods := []jsvm.Module{
		jsmod.NewHTTP(),
		jsmod.NewConsole(c),
		jsmod.NewGlobal("profile", p.cfg),
	}
	if err := jsvm.SetModule(vm, mods); err != nil {
		return err
	}

	str, _ := io.ReadAll(c.Body())
	_, err := vm.RunString(string(str))
	if err != nil {
		return err
	}

	return nil
}
