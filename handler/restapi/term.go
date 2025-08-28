package restapi

import (
	"github.com/coder/websocket"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/business/service"
)

func NewTerm(svc *service.Term) *Term {
	return &Term{svc: svc}
}

type Term struct {
	svc *service.Term
}

func (api *Term) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/ws/pty").GET(api.pty)
	// r.Route("/ws/ssh").GET(api.ssh)
	return nil
}

func (api *Term) pty(c *ship.Context) error {
	req := new(request.TermResize)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	w, r := c.ResponseWriter(), c.Request()
	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer ws.CloseNow()

	if err = api.svc.PTY(ws, req); err != nil {
		reason := err.Error()
		_ = ws.Close(websocket.CloseStatus(err), reason)
	}

	return nil
}

//func (api *Term) ssh(c *ship.Context) error {
//	req := new(request.TermSSH)
//	if err := c.BindQuery(req); err != nil {
//		return err
//	}
//
//	w, r := c.ResponseWriter(), c.Request()
//	ws, err := websocket.Accept(w, r, nil)
//	if err != nil {
//		return err
//	}
//	if err = api.svc.SSH(ws, req); err != nil {
//		reason := err.Error()
//		_ = ws.Close(websocket.CloseStatus(err), reason)
//	}
//
//	return nil
//}
