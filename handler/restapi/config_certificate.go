package restapi

import (
	"archive/zip"
	"mime"
	"net/http"
	"strconv"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/business/service"
	"github.com/xmx/aegis-server/handler/shipx"
)

func NewConfigCertificate(svc service.ConfigCertificate) shipx.Router {
	return &configCertificateAPI{svc: svc}
}

type configCertificateAPI struct {
	svc service.ConfigCertificate
}

func (api *configCertificateAPI) Route(r *ship.RouteGroupBuilder) error {
	r.Route("/config/certificates").GET(api.Page)
	r.Route("/config/certificate/download").GET(api.Download)
	r.Route("/config/certificate/refresh").GET(api.Refresh)
	r.Route("/config/certificate/cond").GET(api.Cond)
	r.Route("/config/certificate").
		POST(api.Create).
		PUT(api.Update).
		DELETE(api.Delete)

	return nil
}

func (api *configCertificateAPI) Cond(c *ship.Context) error {
	ret := api.svc.Cond()
	return c.JSON(http.StatusOK, ret)
}

func (api *configCertificateAPI) Page(c *ship.Context) error {
	req := new(request.PageCond)
	if err := c.BindQuery(req); err != nil {
		return err
	}
	ctx := c.Request().Context()
	ret, err := api.svc.Page(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (api *configCertificateAPI) Create(c *ship.Context) error {
	req := new(request.ConfigCertificateCreate)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return api.svc.Create(ctx, req)
}

func (api *configCertificateAPI) Update(c *ship.Context) error {
	req := new(request.ConfigCertificateUpdate)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return api.svc.Update(ctx, req)
}

func (api *configCertificateAPI) Delete(c *ship.Context) error {
	req := new(request.Int64IDs)
	if err := c.BindQuery(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return api.svc.Delete(ctx, req.ID)
}

func (api *configCertificateAPI) Download(c *ship.Context) error {
	req := new(request.Int64IDs)
	if err := c.BindQuery(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	dats, err := api.svc.Find(ctx, req.ID)
	if err != nil {
		return err
	}
	if len(dats) == 0 {
		return ship.ErrNotFound
	}

	size := len(dats)
	name := "certificate"
	extension := ".zip"
	if size == 1 {
		name = dats[0].CommonName
	}
	filename := name + extension
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

	unique := make(map[string]struct{}, size)
	for _, dat := range dats {
		commonName := dat.CommonName
		if _, ok := unique[commonName]; ok {
			sid := strconv.FormatInt(dat.ID, 10)
			commonName = sid + "-" + commonName
		}
		unique[commonName] = struct{}{}

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
	}

	return err
}

func (api *configCertificateAPI) Refresh(c *ship.Context) error {
	ctx := c.Request().Context()
	return api.svc.Refresh(ctx)
}
