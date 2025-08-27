package restapi

import (
	"log/slog"

	"github.com/xmx/aegis-server/business/service"
	"github.com/xmx/aegis-server/contract/request"
	"github.com/xmx/ship"
)

func NewChannel(svc *service.Channel) *Channel {
	return &Channel{
		svc: svc,
	}
}

type Channel struct {
	svc *service.Channel
}

func (chn *Channel) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/channel").POST(chn.open)

	return nil
}

func (chn *Channel) open(c *ship.Context) error {
	req := new(request.ChannelOpen)
	if err := c.Bind(req); err != nil {
		return err
	}

	w, r := c.Response(), c.Request()
	err := chn.svc.Open(w, r, req)
	if err != nil {
		c.Warn("通道建立失败", slog.Any("error", err))
	}

	return err
}
