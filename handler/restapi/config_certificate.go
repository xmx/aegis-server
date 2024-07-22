package restapi

import (
	"archive/zip"
	"mime"
	"net/http"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/business/service"
	"github.com/xmx/aegis-server/handler/shipx"
)

func ConfigCertificate(svc service.ConfigCertificateService) shipx.Register {
	return &configCertificateAPI{svc: svc}
}

type configCertificateAPI struct {
	svc service.ConfigCertificateService
}

func (cc *configCertificateAPI) Register(r shipx.Router) error {
	auth := r.Auth()
	auth.Route("/config/certificates").GET(cc.List)
	auth.Route("/config/certificate/refresh").GET(cc.Refresh)
	auth.Route("/config/certificate").
		GET(cc.Download).
		POST(cc.Create).
		PUT(cc.Update).
		DELETE(cc.Delete)

	return nil
}

func (cc *configCertificateAPI) List(c *ship.Context) error {
	ctx := c.Request().Context()
	ret, err := cc.svc.List(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (cc *configCertificateAPI) Create(c *ship.Context) error {
	req := new(request.ConfigCertificateCreate)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return cc.svc.Create(ctx, req)
}

func (cc *configCertificateAPI) Update(c *ship.Context) error {
	req := new(request.ConfigCertificateUpdate)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return cc.svc.Update(ctx, req)
}

func (cc *configCertificateAPI) Delete(c *ship.Context) error {
	req := new(request.Int64ID)
	if err := c.BindQuery(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return cc.svc.Delete(ctx, req.ID)
}

func (cc *configCertificateAPI) Download(c *ship.Context) error {
	req := new(request.Int64ID)
	if err := c.BindQuery(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	dat, err := cc.svc.Find(ctx, req.ID)
	if err != nil {
		return err
	}

	extension := ".zip"
	commonName := dat.CommonName
	filename := commonName + extension
	contentType := mime.TypeByExtension(extension)
	if contentType == "" {
		contentType = ship.MIMEOctetStream
	}
	disposition := mime.FormatMediaType("attachment", map[string]string{"filename": filename})
	c.SetRespHeader(ship.HeaderContentType, contentType)
	c.SetRespHeader(ship.HeaderContentDisposition, disposition)
	c.WriteHeader(http.StatusOK)

	zw := zip.NewWriter(c.ResponseWriter())
	//goland:noinspection GoUnhandledErrorResult
	defer zw.Close()
	cw, err := zw.Create(commonName + ".crt")
	if err != nil {
		return err
	}
	if _, err = cw.Write([]byte(dat.PublicKey)); err != nil {
		return err
	}
	cw, err = zw.Create(commonName + ".key")
	if err != nil {
		return err
	}
	_, err = cw.Write([]byte(dat.PrivateKey))

	return err
}

func (cc *configCertificateAPI) Refresh(c *ship.Context) error {
	ctx := c.Request().Context()
	return cc.svc.Refresh(ctx)
}
