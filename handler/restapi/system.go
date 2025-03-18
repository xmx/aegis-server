package restapi

import (
	"github.com/xmx/aegis-server/banner"
	"github.com/xmx/ship"
)

func NewSystem() *System {
	return &System{}
}

type System struct{}

func (sys *System) Route(r *ship.RouteGroupBuilder) error {
	r.Route("/system/banner").GET(sys.banner)
	return nil
}

func (sys *System) banner(c *ship.Context) error {
	_, err := banner.ANSI(c)
	return err
}
