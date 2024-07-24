package restapi

import (
	"net/http"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/argument/errcode"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/argument/response"
	"github.com/xmx/aegis-server/handler/shipx"
)

func Tunnel() shipx.Register {
	return &tunnelAPI{}
}

type tunnelAPI struct{}

func (api *tunnelAPI) Register(rt shipx.Router) error {
	rt.Auth().Route("/tunnel").
		POST(api.Connect)

	return nil
}

func (api *tunnelAPI) Connect(c *ship.Context) error {
	req := new(request.TunnelConnect)
	if err := c.Bind(req); err != nil {
		return err
	}
	w := c.ResponseWriter()
	hijacker, ok := w.(http3.Hijacker)
	if !ok {
		return errcode.ErrConnectionHijack
	}
	hc := hijacker.Connection()
	qc, yes := hc.(quic.Connection)
	if !yes {
		return errcode.ErrConnectionHijack
	}
	_ = qc

	res := &response.TunnelConnect{
		Name: "FAKE",
	}

	return c.JSON(http.StatusAccepted, res)
}
