package restapi

import (
	"net/http"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-common/banner"
	"github.com/xmx/aegis-server/config"
)

func NewSystem(cfg *config.Config) *System {
	return &System{
		cfg: cfg,
	}
}

type System struct {
	cfg *config.Config
}

func (sys *System) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/system/banner").GET(sys.banner)
	r.Route("/system/config").GET(sys.config)
	return nil
}

func (sys *System) banner(c *ship.Context) error {
	_, err := banner.ANSI(c)
	return err
}

func (sys *System) config(c *ship.Context) error {
	return c.JSON(http.StatusOK, sys.cfg)
}
