package restapi

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/xmx/aegis-server/application/manager/request"
	"github.com/xmx/aegis-server/application/manager/service"
)

type Agent struct {
	svc *service.Agent
}

func NewAgent(svc *service.Agent) *Agent {
	return &Agent{
		svc: svc,
	}
}

func (agt *Agent) RegisterRoute(r *echo.Group) error {
	r.GET("/agents", agt.page)
	r.GET("/agent", agt.get)
	r.GET("/agent/records", agt.records)

	return nil
}

func (agt *Agent) page(c *echo.Context) error {
	req := new(request.Pages)
	if err := c.Bind(req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	ret, err := agt.svc.Page(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (agt *Agent) get(c *echo.Context) error {
	req := new(request.ObjectID)
	if err := c.Bind(req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	ret, err := agt.svc.Get(ctx, req.ID.UnsafeID())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (agt *Agent) records(c *echo.Context) error {
	req := new(request.IDPages)
	if err := c.Bind(req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	ret, err := agt.svc.Records(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}
