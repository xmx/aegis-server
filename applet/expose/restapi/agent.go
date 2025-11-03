package restapi

import (
	"net/http"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/applet/expose/request"
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

func (agt *Agent) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/agents").GET(agt.page)

	return nil
}

func (agt *Agent) page(c *ship.Context) error {
	req := new(request.PageKeywords)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	ret, err := agt.svc.Page(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}
