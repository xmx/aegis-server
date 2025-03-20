package restapi

import (
	"net/http"
	"path"
	"strings"

	"github.com/xmx/aegis-server/protocol/webdavfs"
	"github.com/xmx/ship"
)

func NewDAV(basepath, dir string) *DAV {
	const subpath = "/dav"
	basepath = strings.TrimRight(basepath, "/")
	prefix := path.Join(basepath, subpath)
	dav := webdavfs.New(dir)

	return &DAV{
		prefix:  prefix,
		subpath: subpath,
		handler: http.StripPrefix(prefix, dav),
	}
}

type DAV struct {
	prefix  string       // HTTP 路由公共路径
	subpath string       // handler 子路径
	handler http.Handler // WebDAV
}

func (d *DAV) Route(r *ship.RouteGroupBuilder) error {
	methods := []string{
		http.MethodOptions, http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodDelete,
		"LOCK", "UNLOCK", "PROPFIND", "PROPPATCH", "MKCOL", "COPY", "MOVE",
	}
	r.Route(d.subpath).Method(d.access, methods...)
	r.Route(d.subpath+"/*wildcard").Method(d.access, methods...)

	return nil
}

func (d *DAV) access(c *ship.Context) error {
	wildcard := c.Param("wildcard")
	w, r := c.ResponseWriter(), c.Request()
	reqPath := r.URL.Path
	if wildcard == "" && !strings.HasSuffix(reqPath, "/") {
		return c.Redirect(http.StatusPermanentRedirect, reqPath+"/")
	}
	d.handler.ServeHTTP(w, r)

	return nil
}
