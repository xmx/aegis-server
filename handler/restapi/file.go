package restapi

import (
	"mime"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/datalayer/gridfs"
	"github.com/xmx/aegis-server/handler/shipx"
)

func NewFile(dbfs gridfs.FS) shipx.Controller {
	return &fileAPI{dbfs: dbfs}
}

type fileAPI struct {
	dbfs gridfs.FS
}

func (api *fileAPI) Register(rt shipx.Router) error {
	auth := rt.Auth()
	auth.Route("/file").
		PUT(api.Upload).
		GET(api.Download)
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
