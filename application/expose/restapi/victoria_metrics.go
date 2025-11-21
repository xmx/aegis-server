package restapi

import (
	"net/http"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/application/expose/request"
	"github.com/xmx/aegis-server/application/expose/service"
)

func NewVictoriaMetrics(svc *service.VictoriaMetrics) *VictoriaMetrics {
	return &VictoriaMetrics{
		svc: svc,
	}
}

type VictoriaMetrics struct {
	svc *service.VictoriaMetrics
}

func (vm *VictoriaMetrics) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/victoria-metrics").
		GET(vm.list).
		POST(vm.create).
		PUT(vm.update).
		DELETE(vm.delete)

	return nil
}

func (vm *VictoriaMetrics) list(c *ship.Context) error {
	ctx := c.Request().Context()
	dat, _ := vm.svc.List(ctx)

	return c.JSON(http.StatusOK, dat)
}

func (vm *VictoriaMetrics) create(c *ship.Context) error {
	req := new(request.VictoriaMetricsUpsert)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return vm.svc.Create(ctx, req)
}

func (vm *VictoriaMetrics) update(c *ship.Context) error {
	req := new(request.VictoriaMetricsUpsert)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return vm.svc.Update(ctx, req)
}

func (vm *VictoriaMetrics) delete(c *ship.Context) error {
	req := new(request.Names)
	if err := c.BindQuery(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return vm.svc.Delete(ctx, req.Name)
}
