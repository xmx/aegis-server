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

func NewAgentRelease(svc *service.AgentRelease, brok *service.Broker) *AgentRelease {
	return &AgentRelease{
		svc:  svc,
		brok: brok,
	}
}

type AgentRelease struct {
	svc  *service.AgentRelease
	brok *service.Broker
}

func (ar *AgentRelease) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/agent-release").
		GET(ar.download).
		POST(ar.upload).
		DELETE(ar.delete)
	r.Route("/agent-release/parse").POST(ar.parse)

	return nil
}

func (ar *AgentRelease) delete(c *ship.Context) error {
	req := new(request.ObjectID)
	if err := c.BindQuery(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return ar.svc.Delete(ctx, req.OID())
}

func (ar *AgentRelease) upload(c *ship.Context) error {
	req := new(request.AgentReleaseUpload)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return ar.svc.Upload(ctx, req)
}

func (ar *AgentRelease) parse(c *ship.Context) error {
	req := new(request.AgentReleaseUpload)
	if err := c.Bind(req); err != nil {
		return err
	}

	file, err := req.File.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	_, ret, err := ar.svc.Parse(file)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (ar *AgentRelease) download(c *ship.Context) error {
	req := new(request.ReleaseDownload)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	goos, goarch := req.Goos, req.Goarch
	attrs := []any{"goos", goos, "goarch", goarch}
	ctx := c.Request().Context()
	exposes, err := ar.svc.Exposes(ctx)
	if err != nil {
		attrs = append(attrs, "error", err)
		c.Warnf("请检查 broker 是否存在并且配置了暴露地址", attrs...)
		return err
	}

	last, err := ar.svc.Latest(ctx, req.Goos, req.Goarch)
	if err != nil {
		attrs = append(attrs, "error", err)
		c.Warnf("没有找到对应的发行版本", attrs...)
		return err
	}

	stm, err := ar.svc.Open(ctx, last.FileID)
	if err != nil {
		attrs = append(attrs, "error", err)
		c.Warnf("打开文件出错", attrs...)
		return err
	}
	defer stm.Close()

	filesize := last.Length
	manifest := &response.AgentManifest{
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
