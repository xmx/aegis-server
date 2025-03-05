package restapi

import (
	"net/http"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/business/service"
)

func NewOplog(svc *service.Oplog) *Oplog {
	return &Oplog{svc: svc}
}

type Oplog struct {
	svc *service.Oplog
}

func (l *Oplog) Route(r *ship.RouteGroupBuilder) error {
	r.Route("/oplog").GET(l.detail).DELETE(l.delete)
	r.Route("/oplogs").GET(l.page)
	r.Route("/oplog/cond").GET(l.cond)

	return nil
}

func (l *Oplog) cond(c *ship.Context) error {
	ret := l.svc.Cond()
	return c.JSON(http.StatusOK, ret)
}

func (l *Oplog) page(c *ship.Context) error {
	req := new(request.Pages)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	ret, err := l.svc.Page(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (l *Oplog) detail(c *ship.Context) error {
	req := new(request.Int64ID)
	if err := c.BindQuery(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	ret, err := l.svc.Detail(ctx, req.ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (l *Oplog) delete(c *ship.Context) error {
	req := new(request.Pages)
	if err := c.BindQuery(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return l.svc.Delete(ctx, req)
}
