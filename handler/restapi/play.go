package restapi

import (
	"io"
	"os"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/jsenv/jsmod"
	"github.com/xmx/aegis-server/jsenv/jsvm"
)

func NewPlay() *Play {
	return &Play{}
}

type Play struct{}

func (p *Play) Route(r *ship.RouteGroupBuilder) error {
	r.Route("/play/js").POST(p.run)
	return nil
}

func (p *Play) run(c *ship.Context) error {
	vm := jsvm.New()
	jsmod.RegisterHTTP(vm)
	jsmod.RegisterConsole(vm, os.Stdout)

	str, _ := io.ReadAll(c.Body())
	_, err := vm.RunString(string(str))
	if err != nil {
		return err
	}

	return nil
}
