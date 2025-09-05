package restapi

import (
	"net/http"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/profile"
)

func NewConfig(cfg *profile.Config) *Config {
	return &Config{
		cfg: cfg,
	}
}

type Config struct {
	cfg *profile.Config
}

func (cf *Config) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/config").GET(cf.get)
	return nil
}

func (cf *Config) get(c *ship.Context) error {
	ret := cf.cfg.Database
	return c.JSON(http.StatusOK, ret)
}
