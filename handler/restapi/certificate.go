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

func NewCertificate(svc *service.Certificate) *Certificate {
	return &Certificate{svc: svc}
}

type Certificate struct {
	svc *service.Certificate
}

func (crt *Certificate) Route(r *ship.RouteGroupBuilder) error {
	r.Route("/certificates").GET(crt.page)
	r.Route("/certificate/download").GET(crt.download)
	r.Route("/certificate/refresh").GET(crt.refresh)
	r.Route("/certificate/cond").GET(crt.cond)
	r.Route("/certificate").
		GET(crt.detail).
		POST(crt.create).
		PUT(crt.update).
		DELETE(crt.delete)

	return nil
}

func (crt *Certificate) cond(c *ship.Context) error {
	ret := crt.svc.Cond()
	return c.JSON(http.StatusOK, ret)
}

func (crt *Certificate) page(c *ship.Context) error {
	req := new(request.Pages)
	if err := c.BindQuery(req); err != nil {
		return err
	}
	ctx := c.Request().Context()
	ret, err := crt.svc.Page(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (crt *Certificate) create(c *ship.Context) error {
	req := new(request.ConfigCertificateCreate)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return crt.svc.Create(ctx, req)
}

func (crt *Certificate) update(c *ship.Context) error {
	req := new(request.ConfigCertificateUpdate)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return crt.svc.Update(ctx, req)
}

func (crt *Certificate) detail(c *ship.Context) error {
	req := new(request.Int64ID)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	ret, err := crt.svc.Detail(ctx, req.ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (crt *Certificate) delete(c *ship.Context) error {
	req := new(request.Int64IDs)
	if err := c.BindQuery(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return crt.svc.Delete(ctx, req.ID)
}

func (crt *Certificate) download(c *ship.Context) error {
	req := new(request.Int64IDs)
	if err := c.BindQuery(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	dats, err := crt.svc.Find(ctx, req.ID)
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

func (crt *Certificate) refresh(c *ship.Context) error {
	ctx := c.Request().Context()
	_, err := crt.svc.Refresh(ctx)
	return err
}
