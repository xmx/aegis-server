package brkapi

import (
	"net/http"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/business/bservice"
	"github.com/xmx/aegis-server/channel/transport"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func NewAlive(svc *bservice.Alive) *Alive {
	return &Alive{
		svc: svc,
	}
}

type Alive struct {
	svc *bservice.Alive
}

func (alv *Alive) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/alive/ping").GET(alv.ping)
	return nil
}

func (alv *Alive) ping(c *ship.Context) error {
	ctx := c.Request().Context()
	peer := transport.FromContext[bson.ObjectID](ctx)
	_ = alv.svc.Ping(ctx, peer)

	return c.NoContent(http.StatusNoContent)
}
