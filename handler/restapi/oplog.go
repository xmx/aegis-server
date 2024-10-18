package restapi

import (
	"net/http"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/business/service"
	"github.com/xmx/aegis-server/handler/shipx"
)

func NewOplog(svc service.Oplog) shipx.Router {
	return &oplogAPI{svc: svc}
}

type oplogAPI struct {
	svc service.Oplog
}

func (api *oplogAPI) Route(r *ship.RouteGroupBuilder) error {
	r.Route("/oplog").GET(api.Detail).DELETE(api.Delete)
	r.Route("/oplogs").GET(api.Page)
	r.Route("/oplog/cond").GET(api.Cond)

	return nil
}

func (api *oplogAPI) Cond(c *ship.Context) error {
	ret := api.svc.Cond()
	return c.JSON(http.StatusOK, ret)
}

func (api *oplogAPI) Page(c *ship.Context) error {
	req := new(request.PageCondition)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	ret, err := api.svc.Page(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (api *oplogAPI) Detail(c *ship.Context) error {
	req := new(request.Int64ID)
	if err := c.BindQuery(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	ret, err := api.svc.Detail(ctx, req.ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (api *oplogAPI) Delete(c *ship.Context) error {
	req := new(request.CondWhereInputs)
	if err := c.BindQuery(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return api.svc.Delete(ctx, req)
}
