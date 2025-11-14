package restapi

import (
	"io"
	"mime"
	"net/http"
	"strconv"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-common/stegano"
	"github.com/xmx/aegis-server/application/expose/request"
	"github.com/xmx/aegis-server/application/expose/response"
	"github.com/xmx/aegis-server/application/expose/service"
)

func NewBrokerRelease(svc *service.BrokerRelease, brk *service.Broker) *BrokerRelease {
	return &BrokerRelease{
		svc: svc,
		brk: brk,
	}
}

type BrokerRelease struct {
	svc *service.BrokerRelease
	brk *service.Broker
}

func (br *BrokerRelease) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/broker-release").
		GET(br.download).
		POST(br.upload).
		DELETE(br.delete)
	r.Route("/broker-releases").GET(br.list)
	r.Route("/broker-release/parse").POST(br.parse)

	return nil
}

func (br *BrokerRelease) list(c *ship.Context) error {
	ctx := c.Request().Context()
	ret, err := br.svc.List(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (br *BrokerRelease) delete(c *ship.Context) error {
	req := new(request.ObjectID)
	if err := c.BindQuery(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return br.svc.Delete(ctx, req.OID())
}

func (br *BrokerRelease) upload(c *ship.Context) error {
	req := new(request.BrokerReleaseUpload)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return br.svc.Upload(ctx, req)
}

func (br *BrokerRelease) parse(c *ship.Context) error {
	req := new(request.BrokerReleaseUpload)
	if err := c.Bind(req); err != nil {
		return err
	}

	file, err := req.File.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	_, ret, err := br.svc.Parse(file)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (br *BrokerRelease) download(c *ship.Context) error {
	req := new(request.ReleaseBrokerDownload)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	name, goos, goarch := req.Name, req.Goos, req.Goarch
	attrs := []any{"name", name, "goos", goos, "goarch", goarch}
	ctx := c.Request().Context()
	brok, err := br.brk.GetByName(ctx, name)
	if err != nil {
		c.Warnf("请检查 broker 是否注册", attrs...)
		return err
	}
	exposes, err := br.svc.Exposes(ctx)
	if err != nil {
		attrs = append(attrs, "error", err)
		c.Warnf("请检查管理端是否配置了暴露地址", attrs...)
		return err
	}

	last, err := br.svc.Latest(ctx, req.Goos, req.Goarch)
	if err != nil {
		attrs = append(attrs, "error", err)
		c.Warnf("没有找到对应的发行版本", attrs...)
		return err
	}

	stm, err := br.svc.Open(ctx, last.FileID)
	if err != nil {
		attrs = append(attrs, "error", err)
		c.Warnf("打开文件出错", attrs...)
		return err
	}
	defer stm.Close()

	filesize := last.Length
	manifest := &response.BrokerManifest{
		Secret:    brok.Secret,
		Addresses: exposes.Addresses(),
		Offset:    filesize,
	}
	zipbuf, err := stegano.CreateManifestZip(manifest, filesize)
	if err != nil {
		return err
	}

	totalLen := filesize + int64(zipbuf.Len())
	contentLength := strconv.FormatInt(totalLen, 10)
	params := last.Checksum.Map()
	params["filename"] = last.Filename
	mediaType := mime.FormatMediaType("attachment", params)
	c.SetRespHeader(ship.HeaderContentDisposition, mediaType)
	c.SetRespHeader(ship.HeaderContentLength, contentLength)
	down := io.MultiReader(stm, zipbuf)

	return c.Stream(http.StatusOK, ship.MIMEOctetStream, down)
}
