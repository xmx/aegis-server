package restapi

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/argument/errcode"
	"github.com/xmx/aegis-server/handler/shipx"
)

func Transport() shipx.Register {
	return new(transportAPI)
}

type transportAPI struct{}

func (tr *transportAPI) Register(r shipx.Router) error {
	auth := r.Auth()
	auth.Route("/transport").
		GET(tr.Connect).
		CONNECT(tr.Connect)

	return nil
}

func (tr *transportAPI) Connect(c *ship.Context) error {
	w := c.ResponseWriter()
	hc, ok := w.(http3.Hijacker)
	if !ok {
		return errcode.ErrUnsupportedHijack
	}

	data, _ := httputil.DumpRequest(c.Request(), false)
	fmt.Println(string(data))

	qc, yes := hc.Connection().(quic.Connection)
	if !yes {
		return errcode.ErrUnsupportedHijack
	}
	_ = qc
	http.ReadResponse()

	return nil
}
