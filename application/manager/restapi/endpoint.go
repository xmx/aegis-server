package restapi

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/application/muxgate/linkhub"
	"github.com/xmx/muxconn"
)

type Endpoint struct {
	next  linkhub.MUXHandler
	wsupg websocket.Upgrader
}

func NewEndpoint() *Endpoint {
	return &Endpoint{
		wsupg: websocket.Upgrader{
			HandshakeTimeout: 5 * time.Second,
			Subprotocols:     []string{"smux", "yamux"},
		},
	}
}

// connect Agent 接入端点。
//
//goland:noinspection GoUnhandledErrorResult
func (ep *Endpoint) connect(c *ship.Context) error {
	r := c.Request()
	ws, err := ep.wsupg.Upgrade(c, r, nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	subprotocol := ws.Subprotocol()
	if subprotocol == "" {
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

	if err = ep.next.HandleMUX(mux); err != nil {
		c.Warnf("连接已断开：%s", err)
	}

	return nil
}
