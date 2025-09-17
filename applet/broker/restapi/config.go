package restapi

import (
	"net/http"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-control/contract/sbmesg"
	"github.com/xmx/aegis-server/config"
)

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		cfg: cfg,
	}
}

type Config struct {
	cfg *config.Config
}

func (cf *Config) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/config").GET(cf.get)
	return nil
}

func (cf *Config) get(c *ship.Context) error {
	ret := sbmesg.BrokerInitialConfig{
		MongoURI: cf.cfg.Database.URI,
	}

	return c.JSON(http.StatusOK, ret)
}
