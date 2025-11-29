package restapi

import (
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/application/expose/request"
	"github.com/xmx/aegis-server/application/expose/service"
)

func NewFirewall(svc *service.Firewall) *Firewall {
	return &Firewall{
		svc: svc,
	}
}

type Firewall struct {
	svc *service.Firewall
}

func (fw *Firewall) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/log/watch").GET(fw.create)

	return nil
}

func (fw *Firewall) create(c *ship.Context) error {
	req := new(request.FirewallUpsert)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return fw.svc.Create(ctx, req)
}

func (fw *Firewall) update(c *ship.Context) error {
	req := new(request.FirewallUpsert)
	if err := c.Bind(req); err != nil {
		return err
	}
	return nil
}

func (fw *Firewall) delete(c *ship.Context) error {
	req := new(request.Names)
	if err := c.BindQuery(req); err != nil {
		return err
	}
	return nil
}
