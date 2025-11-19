package restapi

import (
	"github.com/gorilla/websocket"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-common/library/httpkit"
	"github.com/xmx/aegis-common/tunnel/tunconst"
	"github.com/xmx/aegis-common/tunnel/tunopen"
)

func NewTunnel(next tunconst.Handler) *Tunnel {
	return &Tunnel{
		next: next,
		wsup: httpkit.NewWebsocketUpgrader(),
	}
}

type Tunnel struct {
	next tunconst.Handler
	wsup *websocket.Upgrader
}

func (tnl *Tunnel) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/tunnel").GET(tnl.open)

	return nil
}

func (tnl *Tunnel) open(c *ship.Context) error {
	w, r := c.ResponseWriter(), c.Request()
	ws, err := tnl.wsup.Upgrade(w, r, nil)
	if err != nil {
		c.Warnf("websocket 升级失败", "error", err)
		return nil // 无需再次返回 err 信息
	}
	conn := ws.NetConn()
	mux, err := tunopen.NewSMUX(conn, nil, true)
	if err != nil {
		_ = conn.Close()
		return err
	}
	tnl.next.Handle(mux)

	return nil
}
