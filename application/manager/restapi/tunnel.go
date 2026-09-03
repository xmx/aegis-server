package restapi

import (
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
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

func (tnl *Tunnel) RegisterRoute(r *echo.Group) error {
	r.GET("/tunnel", tnl.connect)

	return nil
}

// connect Agent 接入端点。
//
//goland:noinspection GoUnhandledErrorResult
func (tnl *Tunnel) connect(c *echo.Context) error {
	// clientIP := c.RealIP()

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

	_ = tnl.next.HandleMUX(mux)

	return nil
}
