package restapi

import (
	"mime"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-control/datalayer/model"
	"github.com/xmx/aegis-server/application/expose/request"
	"github.com/xmx/aegis-server/application/expose/service"
)

func NewFS(svc *service.FS) *FS {
	return &FS{
		svc: svc,
	}
}

type FS struct {
	svc *service.FS
}

func (fs *FS) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/fs/http").GET(fs.http)
	r.Route("/fs/http/*path").GET(fs.http)
	r.Route("/fs/list").GET(fs.list)
	r.Route("/fs/list/*path").GET(fs.list)
	r.Route("/fs/download").GET(fs.download)
	r.Route("/fs/download/*path").GET(fs.download)
	r.Route("/fs/create").PUT(fs.create)
	r.Route("/fs/create/*path").PUT(fs.create)
	r.Route("/fs/update").PUT(fs.update)
	r.Route("/fs/update/*path").PUT(fs.update)
	r.Route("/fs/mkdir").POST(fs.mkdir)
	r.Route("/fs/mkdir/*path").POST(fs.mkdir)
	r.Route("/fs/remove").DELETE(fs.remove)
	r.Route("/fs/remove/*path").DELETE(fs.remove)

	return nil
}

func (fs *FS) create(c *ship.Context) error {
	req := new(request.FSUpload)
	if err := c.Bind(req); err != nil {
		return err
	}
	dir := c.Param("path")
	ctx := c.Request().Context()

	ret, err := fs.svc.Create(ctx, dir, req.File)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (fs *FS) update(c *ship.Context) error {
	req := new(request.FSUpload)
	if err := c.Bind(req); err != nil {
		return err
	}
	dir := c.Param("path")
	ctx := c.Request().Context()

	ret, err := fs.svc.Update(ctx, dir, req.File)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (fs *FS) mkdir(c *ship.Context) error {
	dir := c.Param("path")
	ctx := c.Request().Context()
	err := fs.svc.Mkdir(ctx, dir)

	return err
}

func (fs *FS) list(c *ship.Context) error {
	dir := c.Param("path")
	ctx := c.Request().Context()
	ret, err := fs.svc.List(ctx, dir)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (fs *FS) download(c *ship.Context) error {
	dir := c.Param("path")
	ctx := c.Request().Context()
	f, err := fs.svc.Open(ctx, dir)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer f.Close()

	var contentType string
	if stat, _ := f.Stat(); stat != nil {
		name := stat.Name()
		ext := strings.ToLower(path.Ext(name))
		contentType = mime.TypeByExtension(ext)

		size := strconv.FormatInt(stat.Size(), 10)
		c.SetRespHeader(ship.HeaderContentLength, size)

		params := make(map[string]string, 4)
		if mfs, _ := stat.Sys().(*model.FS); mfs != nil {
			params = mfs.Checksum.Map()
		}
		params["filename"] = name
		disposition := mime.FormatMediaType("attachment", params)
		c.SetRespHeader(ship.HeaderContentDisposition, disposition)
	}
	if contentType == "" {
		contentType = ship.MIMEOctetStream
	}

	return c.Stream(http.StatusOK, contentType, f)
}

func (fs *FS) remove(c *ship.Context) error {
	dir := c.Param("path")
	ctx := c.Request().Context()

	return fs.svc.Remove(ctx, dir)
}

func (fs *FS) http(c *ship.Context) error {
	dir := "/" + c.Param("path")
	w, r := c.Response(), c.Request()
	r.URL.Path = dir
	fs.svc.Handler().ServeHTTP(w, r)
	return nil
}
