package restapi

import (
	"net/http"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/applet/expose/service"
	"github.com/xmx/aegis-server/contract/request"
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
	r.Route("/brokers").GET(bk.list)
	r.Route("/broker").POST(bk.create)
	r.Route("/broker/kickout").GET(bk.kickout)

	return nil
}

func (bk *Broker) create(c *ship.Context) error {
	req := new(request.BrokerCreate)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return bk.svc.Create(ctx, req.Name)
}

func (bk *Broker) list(c *ship.Context) error {
	ctx := c.Request().Context()
	ret, err := bk.svc.List(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (bk *Broker) kickout(c *ship.Context) error {
	req := new(request.ObjectID)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	return bk.svc.Kickout(req.OID())
}
