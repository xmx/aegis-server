package restapi

import (
	"net/http"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/applet/expose/request"
	"github.com/xmx/aegis-server/applet/expose/service"
)

func NewBroker(svc *service.Broker) *Broker {
	return &Broker{
		svc: svc,
	}
}

type Broker struct {
	svc *service.Broker
}

func (bk *Broker) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/brokers").GET(bk.page)
	r.Route("/broker").POST(bk.create)

	return nil
}

func (bk *Broker) page(c *ship.Context) error {
	req := new(request.PageKeywords)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	ret, err := bk.svc.Page(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (bk *Broker) create(c *ship.Context) error {
	req := new(request.BrokerCreate)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return bk.svc.Create(ctx, req.Name)
}

func (bk *Broker) kickout(c *ship.Context) error {
	req := new(request.ObjectID)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	return bk.svc.Kickout(req.OID())
}
