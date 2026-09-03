package restapi

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/xmx/aegis-server/application/manager/request"
	"github.com/xmx/aegis-server/application/manager/service"
)

type OTelClient struct {
	svc *service.OTelClient
}

func NewOTelClient(svc *service.OTelClient) *OTelClient {
	return &OTelClient{
		svc: svc,
	}
}

func (oc *OTelClient) RegisterRoute(r *echo.Group) error {
	r.GET("/otel-clients", oc.list)
	r.POST("/otel-client", oc.create)
	r.PUT("/otel-client", oc.update)
	r.DELETE("/otel-client/:id", oc.delete)

	return nil
}

func (oc *OTelClient) list(c *echo.Context) error {
	ctx := c.Request().Context()
	ret, err := oc.svc.List(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (oc *OTelClient) create(c *echo.Context) error {
	req := new(request.OTelClientCreate)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return oc.svc.Create(ctx, req)
}

func (oc *OTelClient) update(c *echo.Context) error {
	req := new(request.OTelClientUpdate)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return oc.svc.Update(ctx, req)
}

func (oc *OTelClient) delete(c *echo.Context) error {
	req := new(request.ObjectID)
	if err := c.Bind(req); err != nil {
		return err
	}

	ctx := c.Request().Context()

	return oc.svc.Delete(ctx, req.ID.UnsafeID())
}
