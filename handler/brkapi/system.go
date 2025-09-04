package brkapi

import (
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-common/transport"
	"github.com/xmx/aegis-server/business/bservice"
	"github.com/xmx/aegis-server/contract/brequest"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func NewSystem(svc *bservice.System) *System {
	return &System{
		svc: svc,
	}
}

type System struct {
	svc *bservice.System
}

func (stm *System) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/network-card").POST(stm.networkCard)
	return nil
}

func (stm *System) networkCard(c *ship.Context) error {
	req := new(brequest.SystemNetworkCard)
	if err := c.Bind(req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	peer := transport.FromContext[bson.ObjectID](ctx)

	return stm.svc.NetworkCard(ctx, req, peer)
}
