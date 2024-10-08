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

func NewFile(dbfs gridfs.FS, svc service.File) shipx.Controller {
	return &fileAPI{
		dbfs: dbfs,
		svc:  svc,
	}
}

type fileAPI struct {
	dbfs gridfs.FS
	svc  service.File
}

func (api *fileAPI) Register(rt shipx.Router) error {
	auth := rt.Auth()
	auth.Route("/file").
		Name("文件管理").
		Data(map[string]string{"key": "上传下载"}).
		PUT(api.Upload).
		GET(api.Download)
	auth.Route("/file/cond").GET(api.Cond)
	auth.Route("/files").
		GET(api.Page)

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
	param := map[string]string{
		"filename": filename,
		"sha1":     file.SHA1(),
		"sha256":   file.SHA256(),
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
	req := new(request.PageKeywordOrder)
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
