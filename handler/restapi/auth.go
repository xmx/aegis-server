package restapi

import (
	"log/slog"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/handler/shipx"
)

func NewAuth() shipx.Router {
	return &authAPI{}
}

type authAPI struct{}

func (api *authAPI) Route(r *ship.RouteGroupBuilder) error {
	r.Route("/auth/back").GET(api.Back)
	return nil
}

func (api *authAPI) Back(c *ship.Context) error {
	code := c.Query("code")
	c.Infof("GitHub code", slog.String("code", code))
	return nil
}
