package restapi

import (
	"archive/zip"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/xmx/aegis-server/business/service"
	"github.com/xmx/aegis-server/contract/request"
	"github.com/xmx/ship"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func NewCertificate(svc *service.Certificate) *Certificate {
	return &Certificate{svc: svc}
}

type Certificate struct {
	svc *service.Certificate
}

func (crt *Certificate) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/certificates").GET(crt.page)
	r.Route("/certificate/download").GET(crt.download)
	r.Route("/certificate/forget").DELETE(crt.forget)
	r.Route("/certificate/parse").POST(crt.parse)
	r.Route("/certificate").
		GET(crt.detail).
		POST(crt.create).
		PUT(crt.update).
		DELETE(crt.delete)

	return nil
}

func (crt *Certificate) page(c *ship.Context) error {
	req := new(request.PageKeywords)
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

func (crt *Certificate) parse(c *ship.Context) error {
	req := new(request.ConfigCertificateParse)
	if err := c.Bind(req); err != nil {
		return err
	}

	ret, err := crt.svc.Parse(req.PublicKey, req.PrivateKey)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
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
	req := new(request.ObjectID)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	oid, _ := bson.ObjectIDFromHex(req.ID)
	ret, err := crt.svc.Detail(ctx, oid)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (crt *Certificate) delete(c *ship.Context) error {
	req := new(request.ObjectIDs)
	if err := c.BindQuery(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return crt.svc.Delete(ctx, req.OIDs())
}

func (crt *Certificate) download(c *ship.Context) error {
	req := new(request.ObjectIDs)
	if err := c.BindQuery(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	now := time.Now()
	name := "certificate-" + strconv.FormatInt(now.Unix(), 10)
	extension := ".zip"
	filename := name + extension
	contentType := mime.TypeByExtension(extension)
	if contentType == "" {
		contentType = ship.MIMEOctetStream
	}
	disposition := mime.FormatMediaType("attachment", map[string]string{"filename": filename})
	c.SetRespHeader(ship.HeaderContentType, contentType)
	c.SetRespHeader(ship.HeaderContentDisposition, disposition)
	c.WriteHeader(http.StatusOK)

	zipfs := zip.NewWriter(c.ResponseWriter())
	//goland:noinspection GoUnhandledErrorResult
	defer zipfs.Close()

	// 这些字符不允许当作文件名。
	replacer := strings.NewReplacer("\\", "_", "/", "_", ":", "_", "*", "_",
		"?", "_", "\"", "_", "<", "_", ">", "_", "|", "_")

	names := make(map[string]struct{}, 16)
	for dat, exx := range crt.svc.All(ctx, req.OIDs()) {
		if exx != nil {
			return exx
		}

		cname := replacer.Replace(dat.Name)
		if _, ok := names[cname]; ok {
			sid := dat.ID.Hex()
			cname = cname + "-" + sid
		}
		names[cname] = struct{}{}

		cw, err := zipfs.Create(cname + ".crt")
		if err != nil {
			return err
		}

		pubKey, priKey := []byte(dat.PublicKey), []byte(dat.PrivateKey)
		if _, err = cw.Write(pubKey); err != nil {
			return err
		}
		cw, err = zipfs.Create(cname + ".key")
		if err != nil {
			return err
		}
		_, err = cw.Write(priKey)
	}

	return nil
}

func (crt *Certificate) forget(c *ship.Context) error {
	ctx := c.Request().Context()
	crt.svc.Forget(ctx)
	return nil
}
