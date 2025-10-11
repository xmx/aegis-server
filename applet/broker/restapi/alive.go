package restapi

import (
	"net/http"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-control/linkhub"
	"github.com/xmx/aegis-server/business/bservice"
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
	peer, _ := linkhub.FromContext(ctx)
	_ = alv.svc.Ping(ctx, peer)

	return c.NoContent(http.StatusNoContent)
}
