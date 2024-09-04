package restapi

import (
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/handler/shipx"
)

func NewVNC() shipx.Controller {
	return &vncAPI{}
}

type vncAPI struct{}

func (api *vncAPI) Register(rt shipx.Router) error {
	auth := rt.Auth()
	auth.Route("/ws/vnc").GET(api.VNC)
	return nil
}

func (api *vncAPI) VNC(c *ship.Context) error {
	return nil
}
