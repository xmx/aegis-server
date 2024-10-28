package restapi

import (
	"mime"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/business/service"
	"github.com/xmx/aegis-server/datalayer/gridfs"
	"github.com/xmx/aegis-server/handler/shipx"
)

func NewFile(dbfs gridfs.FS, svc *service.File) shipx.Router {
	return &fileAPI{
		dbfs: dbfs,
		svc:  svc,
	}
}

type fileAPI struct {
	dbfs gridfs.FS
	svc  *service.File
}

func (api *fileAPI) Route(r *ship.RouteGroupBuilder) error {
	r.Route("/file").
		Name("文件管理").
		Data(map[string]string{"key": "上传下载"}).
		PUT(api.Upload).
		GET(api.Download)
	r.Route("/file/cond").GET(api.Cond)
	r.Route("/files").GET(api.Page)
	r.Route("/file/count").GET(api.Count)

	return nil
}

func (api *fileAPI) Upload(c *ship.Context) error {
	file, header, err := c.FormFile("file")
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	ctx := c.Request().Context()
	filename := header.Filename
	_, err = api.dbfs.Save(ctx, filename, file)

	return err
}

func (api *fileAPI) Download(c *ship.Context) error {
	file, err := api.dbfs.OpenID(1)
	if err != nil {
		return err
	}
	defer file.Close()

	filename := file.Name()
	param := map[string]string{"filename": filename}
	if digest, ok := file.(gridfs.Digester); ok {
		param["md5"] = digest.MD5()
		param["sha1"] = digest.SHA1()
		param["sha256"] = digest.SHA256()
	}

	disposition := mime.FormatMediaType("attachment", param)
	c.SetRespHeader(ship.HeaderContentDisposition, disposition)
	c.SetRespHeader(ship.HeaderContentLength, strconv.FormatInt(file.Size(), 10))
	ext := filepath.Ext(filename)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = ship.MIMEOctetStream
	}

	return c.Stream(http.StatusOK, contentType, file)
}

func (api *fileAPI) Page(c *ship.Context) error {
	req := new(request.PageCondition)
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

func (api *fileAPI) Cond(c *ship.Context) error {
	ret := api.svc.Cond()
	return c.JSON(http.StatusOK, ret)
}

func (api *fileAPI) Count(c *ship.Context) error {
	ctx := c.Request().Context()
	ret, err := api.svc.Count(ctx, 10)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ret)
}
