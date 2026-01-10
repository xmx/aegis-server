package restapi

import (
	"net/http"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/application/expose/request"
	"github.com/xmx/aegis-server/application/expose/service"
)

func NewAgentConnectHistory(svc *service.AgentConnectHistory) *AgentConnectHistory {
	return &AgentConnectHistory{
		svc: svc,
	}
}

type AgentConnectHistory struct {
	svc *service.AgentConnectHistory
}

func (ach *AgentConnectHistory) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/agent/connect/histories").GET(ach.page)

	return nil
}

func (ach *AgentConnectHistory) page(c *ship.Context) error {
	req := new(request.PageKeywords)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	ret, err := ach.svc.Page(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}
