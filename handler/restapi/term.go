package restapi

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/business/service"
)

func NewTerm(svc *service.Term) *Term {
	return &Term{
		svc: svc,
		upg: &websocket.Upgrader{
			HandshakeTimeout:  30 * time.Second,
			CheckOrigin:       func(*http.Request) bool { return true },
			EnableCompression: true,
		},
	}
}

type Term struct {
	svc *service.Term
	upg *websocket.Upgrader
}

func (tm *Term) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/pty").GET(tm.pty)
	return nil
}

func (tm *Term) pty(c *ship.Context) error {
	req := new(request.TermResize)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	w, r := c.ResponseWriter(), c.Request()
	ws, err := tm.upg.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer ws.Close()

	if err = tm.svc.PTY(ws, req); err != nil {
		reason := []byte(err.Error())
		deadline := time.Now().Add(10 * time.Second)
		_ = ws.WriteControl(websocket.CloseNormalClosure, reason, deadline)
	}

	return nil
}
