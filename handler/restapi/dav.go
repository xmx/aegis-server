package restapi

import (
	"net/http"
	"strings"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/protocol/webfs"
)

func NewDAV(prefix, dir string) *DAV {
	prefix = strings.TrimRight(prefix, "/")
	const path = "/dav"
	prefix += path
	dav := webfs.DAV(prefix, dir)

	return &DAV{
		path:   path,
		prefix: prefix,
		dav:    dav,
	}
}

type DAV struct {
	path   string
	prefix string // HTTP path 前缀
	dav    http.Handler
}

func (d *DAV) Route(r *ship.RouteGroupBuilder) error {
	methods := []string{
		http.MethodOptions, http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodDelete,
		"LOCK", "UNLOCK", "PROPFIND", "PROPPATCH", "MKCOL", "COPY", "MOVE",
	}
	r.Route(d.path).Method(d.browser, methods...)
	r.Route(d.path+"/*wildcard").Method(d.browser, methods...)

	return nil
}

func (d *DAV) browser(c *ship.Context) error {
	wildcard := c.Param("path")
	w, r := c.ResponseWriter(), c.Request()
	path := r.URL.Path
	if wildcard == "" && !strings.HasSuffix(path, "/") {
		return c.Redirect(http.StatusPermanentRedirect, path+"/")
	}

	d.dav.ServeHTTP(w, r)

	return nil
}
