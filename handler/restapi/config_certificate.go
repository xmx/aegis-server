package restapi

import (
	"archive/zip"
	"mime"
	"net/http"
	"strconv"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/business/service"
)

func NewConfigCertificate(svc *service.ConfigCertificate) *ConfigCertificate {
	return &ConfigCertificate{svc: svc}
}

type ConfigCertificate struct {
	svc *service.ConfigCertificate
}

func (cc *ConfigCertificate) Route(r *ship.RouteGroupBuilder) error {
	r.Route("/config/certificates").GET(cc.page)
	r.Route("/config/certificate/download").GET(cc.download)
	r.Route("/config/certificate/refresh").GET(cc.refresh)
	r.Route("/config/certificate/cond").GET(cc.cond)
	r.Route("/config/certificate").
		GET(cc.detail).
		POST(cc.create).
		PUT(cc.update).
		DELETE(cc.delete)

	return nil
}

func (cc *ConfigCertificate) cond(c *ship.Context) error {
	ret := cc.svc.Cond()
	return c.JSON(http.StatusOK, ret)
}

func (cc *ConfigCertificate) page(c *ship.Context) error {
	req := new(request.PageCondition)
	if err := c.BindQuery(req); err != nil {
		return err
	}
	ctx := c.Request().Context()
	ret, err := cc.svc.Page(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (cc *ConfigCertificate) create(c *ship.Context) error {
	req := new(request.ConfigCertificateCreate)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return cc.svc.Create(ctx, req)
}

func (cc *ConfigCertificate) update(c *ship.Context) error {
	req := new(request.ConfigCertificateUpdate)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return cc.svc.Update(ctx, req)
}

func (cc *ConfigCertificate) detail(c *ship.Context) error {
	req := new(request.Int64ID)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	ret, err := cc.svc.Detail(ctx, req.ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (cc *ConfigCertificate) delete(c *ship.Context) error {
	req := new(request.Int64IDs)
	if err := c.BindQuery(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return cc.svc.Delete(ctx, req.ID)
}

func (cc *ConfigCertificate) download(c *ship.Context) error {
	req := new(request.Int64IDs)
	if err := c.BindQuery(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	dats, err := cc.svc.Find(ctx, req.ID)
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

func (cc *ConfigCertificate) refresh(c *ship.Context) error {
	ctx := c.Request().Context()
	_, err := cc.svc.Refresh(ctx)
	return err
}
