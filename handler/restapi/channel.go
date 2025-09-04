package restapi

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-common/transport"
)

func NewChannel(next transport.Handler) *Channel {
	return &Channel{
		next: next,
		upg: &websocket.Upgrader{
			HandshakeTimeout: 10 * time.Second,
			CheckOrigin:      func(*http.Request) bool { return true },
		},
	}
}

type Channel struct {
	next transport.Handler
	upg  *websocket.Upgrader
}

func (chn *Channel) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/channel").GET(chn.open)

	return nil
}

func (chn *Channel) open(c *ship.Context) error {
	w, r := c.ResponseWriter(), c.Request()
	ws, err := chn.upg.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	conn := ws.NetConn()
	mux, err := transport.NewSMUX(conn, true)
	if err != nil {
		_ = conn.Close()
		return err
	}

	if err = chn.next.Handle(mux); err != nil {
		c.Infof("处理发生错误", slog.Any("error", err))
	}

	return nil
}
