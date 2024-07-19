package restapi

import (
	"github.com/xgfone/ship/v5"
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
	return nil
}
