package restapi

import (
	"net/http"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/application/expose/request"
	"github.com/xmx/aegis-server/application/expose/service"
)

func NewPyroscope(svc *service.Pyroscope) *Pyroscope {
	return &Pyroscope{
		svc: svc,
	}
}

type Pyroscope struct {
	svc *service.Pyroscope
}

func (prs *Pyroscope) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/pyroscope").
		GET(prs.list).
		POST(prs.create).
		PUT(prs.update).
		DELETE(prs.delete)

	return nil
}

func (prs *Pyroscope) list(c *ship.Context) error {
	ctx := c.Request().Context()
	dat, _ := prs.svc.List(ctx)

	return c.JSON(http.StatusOK, dat)
}

func (prs *Pyroscope) create(c *ship.Context) error {
	req := new(request.PyroscopeUpsert)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return prs.svc.Create(ctx, req)
}

func (prs *Pyroscope) update(c *ship.Context) error {
	req := new(request.PyroscopeUpsert)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return prs.svc.Update(ctx, req)
}

func (prs *Pyroscope) delete(c *ship.Context) error {
	req := new(request.Names)
	if err := c.BindQuery(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return prs.svc.Delete(ctx, req.Name)
}
