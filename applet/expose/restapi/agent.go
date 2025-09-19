package restapi

import (
	"net/http"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/applet/expose/service"
)

func NewAgent(svc *service.Agent) *Agent {
	return &Agent{
		svc: svc,
	}
}

type Agent struct {
	svc *service.Agent
}

func (a *Agent) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/agents").GET(a.list)

	return nil
}

func (a *Agent) list(c *ship.Context) error {
	ctx := c.Request().Context()
	ret, err := a.svc.List(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}
