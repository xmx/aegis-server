package restapi

import (
	"github.com/gorilla/websocket"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-common/library/httpkit"
	"github.com/xmx/aegis-common/muxlink/muxconn"
	"github.com/xmx/aegis-common/muxlink/muxproto"
	"github.com/xmx/aegis-common/shipx"
)

func NewTunnel(acpt muxproto.MUXAccepter) *Tunnel {
	return &Tunnel{
		acpt: acpt,
		wsup: httpkit.NewWebsocketUpgrader(),
	}
}

type Tunnel struct {
	acpt muxproto.MUXAccepter
	wsup *websocket.Upgrader
}

func (tnl *Tunnel) RegisterRoute(r *ship.RouteGroupBuilder) error {
	data := shipx.NewRouteData("通道接入点（wss）")
	r.Route("/tunnel").Data(data).GET(tnl.open)

	return nil
}

//goland:noinspection GoUnhandledErrorResult
func (tnl *Tunnel) open(c *ship.Context) error {
	w, r := c.ResponseWriter(), c.Request()
	ws, err := tnl.wsup.Upgrade(w, r, nil)
	if err != nil {
		c.Warnf("websocket 升级失败", "error", err)
		return nil // 无需再次返回 err 信息
	}

	conn := ws.NetConn()
	mux, err := muxconn.NewSMUX(conn, nil, true)
	if err != nil {
		_ = conn.Close()
		return err
	}
	tnl.acpt.AcceptMUX(mux)

	return nil
}
