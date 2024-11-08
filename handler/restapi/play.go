package restapi

import (
	"github.com/grafana/sobek"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/jsenv/jsrt"
)

func NewPlay() *Play {
	vm := jsrt.New()
	return &Play{
		vm: vm,
	}
}

type Play struct {
	vm *sobek.Runtime
}

func (p *Play) run(c *ship.Context) error {
	return nil
}
