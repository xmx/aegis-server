package restapi

import (
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/handler/shipx"
)

type configCertificateAPI struct{}

func (cc *configCertificateAPI) Register(r shipx.Router) error {
	r.Anon().Route("/config/certificate").POST(cc.Create)
	return nil
}

func (cc *configCertificateAPI) Create(c *ship.Context) error {
	req := new(request.ConfigCertificateCreate)
	if err := c.Bind(req); err != nil {
		return err
	}

	return nil
}
