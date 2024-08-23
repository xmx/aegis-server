package restapi

import (
	"github.com/coder/websocket"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/business/service"
	"github.com/xmx/aegis-server/handler/shipx"
)

func NewTerm(svc service.Term) shipx.Controller {
	return &termAPI{svc: svc}
}

type termAPI struct {
	svc service.Term
}

func (api *termAPI) Register(rt shipx.Router) error {
	auth := rt.Auth()
	auth.Route("/ws/pty").GET(api.PTY)
	auth.Route("/ws/ssh").GET(api.SSH)
	return nil
}

func (api *termAPI) PTY(c *ship.Context) error {
	req := new(request.TermResize)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	w, r := c.ResponseWriter(), c.Request()
	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
	}
	_ = api.svc.PTY(ws, req)

	return nil
}

func (api *termAPI) SSH(c *ship.Context) error {
	req := new(request.TermSSH)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	w, r := c.ResponseWriter(), c.Request()
	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
	}
	_ = api.svc.SSH(ws, req)

	return nil
}
