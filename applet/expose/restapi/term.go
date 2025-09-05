package restapi

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-common/library/wsocket"
	"github.com/xmx/aegis-server/applet/expose/service"
	"github.com/xmx/aegis-server/argument/request"
)

func NewTerm(svc *service.Term) *Term {
	return &Term{
		svc:  svc,
		wsup: wsocket.NewUpgrade(),
	}
}

type Term struct {
	svc  *service.Term
	wsup *websocket.Upgrader
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
	ws, err := tm.wsup.Upgrade(w, r, nil)
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
