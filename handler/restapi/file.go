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
)

func NewFile(svc *service.File) *File {
	return &File{
		svc: svc,
	}
}

type File struct {
	svc *service.File
}

func (f *File) Route(r *ship.RouteGroupBuilder) error {
	r.Route("/file").
		PUT(f.upload).
		GET(f.download)
	r.Route("/file/cond").GET(f.cond)
	r.Route("/files").GET(f.page)
	r.Route("/file/count").GET(f.count)

	return nil
}

func (f *File) upload(c *ship.Context) error {
	file, header, err := c.FormFile("file")
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	ctx := c.Request().Context()
	filename := header.Filename
	_, err = f.svc.Save(ctx, filename, file)

	return err
}

func (f *File) download(c *ship.Context) error {
	req := new(request.Int64ID)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	file, err := f.svc.Open(ctx, req.ID)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
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

func (f *File) page(c *ship.Context) error {
	req := new(request.PageCondition)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	ret, err := f.svc.Page(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (f *File) cond(c *ship.Context) error {
	ret := f.svc.Cond()
	return c.JSON(http.StatusOK, ret)
}

func (f *File) count(c *ship.Context) error {
	ctx := c.Request().Context()
	ret, err := f.svc.Count(ctx, 10)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ret)
}
