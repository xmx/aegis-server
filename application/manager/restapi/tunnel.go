package restapi

import (
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/application/nodelink/linkhub"
	"github.com/xmx/muxconn"
)

type Tunnel struct {
	next  linkhub.MUXHandler
	log   *slog.Logger
	wsupg websocket.Upgrader
}

func NewTunnel(next linkhub.MUXHandler, log *slog.Logger) *Tunnel {
	return &Tunnel{
		next: next,
		log:  log,
		wsupg: websocket.Upgrader{
			HandshakeTimeout: 5 * time.Second,
			Subprotocols:     []string{"smux", "yamux"},
		},
	}
}

func (tnl *Tunnel) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/tunnel").GET(tnl.connect)

	return nil
}

// connect Agent 接入端点。
//
//goland:noinspection GoUnhandledErrorResult
func (tnl *Tunnel) connect(c *ship.Context) error {
	clientIP := c.ClientIP()
	c.Infof("通道准备建连 [%s]", clientIP)

	w, r := c.Response(), c.Request()
	ws, err := tnl.wsupg.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	subprotocol := ws.Subprotocol()
	if subprotocol == "" {
		msg := websocket.FormatCloseMessage(websocket.CloseProtocolError, "subprotocol unmatched")
		deadline := time.Now().Add(3 * time.Second)
		_ = ws.WriteControl(websocket.CloseMessage, msg, deadline)
		return nil
	}

	ctx := r.Context()
	conn := ws.NetConn()
	var mux muxconn.Muxer
	if subprotocol == "yamux" {
		mux, err = muxconn.NewYaMUX(ctx, conn, nil, true)
	} else {
		mux, err = muxconn.NewSMUX(ctx, conn, nil, true)
	}
	if err != nil {
		return err
	}
	defer mux.Close()

	err = tnl.next.HandleMUX(mux)
	c.Warnf("连接已断开 [%s] %v", clientIP, err)

	return nil
}
